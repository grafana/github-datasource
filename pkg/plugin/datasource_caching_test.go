package plugin

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
)

// mockFramer is a struct implementing the Framer interface that returns predefined frames for testing purposes
type mockFramer struct {
	frames data.Frames
}

func (m mockFramer) Frames() data.Frames {
	return m.frames
}

// Fixture for the test cases
var dataQueryA = backend.DataQuery{JSON: json.RawMessage(`{"query": "A"}`)}
var framesA = data.Frames{data.NewFrame("A", nil)}
var dataQueryB = backend.DataQuery{JSON: json.RawMessage(`{"query": "B"}`)}
var framesB = data.Frames{data.NewFrame("B", nil)}

type blockingDeploymentsDatasource struct {
	Datasource
	calls        atomic.Int32
	started      chan struct{}
	release      <-chan struct{}
	checkContext bool
	once         sync.Once
}

func (d *blockingDeploymentsDatasource) HandleDeploymentsQuery(ctx context.Context, _ *models.DeploymentsQuery, _ backend.DataQuery) (dfutil.Framer, error) {
	d.calls.Add(1)
	d.once.Do(func() { close(d.started) })
	<-d.release
	if d.checkContext && ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return mockFramer{frames: framesA}, nil
}

func TestCachedDeploymentsQueryCoalescesCacheMisses(t *testing.T) {
	const queries = 10
	release := make(chan struct{})
	datasource := &blockingDeploymentsDatasource{started: make(chan struct{}), release: release}
	cached := WithCaching(datasource)
	var wg sync.WaitGroup
	start := make(chan struct{})
	errs := make(chan error, queries)

	for range queries {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := cached.HandleDeploymentsQuery(context.Background(), &models.DeploymentsQuery{}, dataQueryA)
			errs <- err
		}()
	}
	close(start)
	<-datasource.started
	time.Sleep(10 * time.Millisecond)
	if calls := datasource.calls.Load(); calls != 1 {
		t.Fatalf("expected one underlying deployment query, got %d", calls)
	}
	close(release)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCachedDeploymentsQueryIgnoresLeaderCancellation(t *testing.T) {
	release := make(chan struct{})
	datasource := &blockingDeploymentsDatasource{started: make(chan struct{}), release: release, checkContext: true}
	cached := WithCaching(datasource)
	leaderCtx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 2)

	go func() {
		_, err := cached.HandleDeploymentsQuery(leaderCtx, &models.DeploymentsQuery{}, dataQueryA)
		errs <- err
	}()
	<-datasource.started
	cancel()
	go func() {
		_, err := cached.HandleDeploymentsQuery(context.Background(), &models.DeploymentsQuery{}, dataQueryA)
		errs <- err
	}()
	close(release)

	for range 2 {
		if err := <-errs; err != nil {
			t.Fatalf("expected shared deployment query to survive leader cancellation: %v", err)
		}
	}
}

func TestWithCaching(t *testing.T) {
	cachedDS := WithCaching(nil)

	t.Run("read from empty cache concurrently", func(t *testing.T) {
		var wg sync.WaitGroup

		// Read goroutine 1
		wg.Add(1)
		go func() {
			defer wg.Done()

			f, err := cachedDS.getCache(dataQueryA)
			assert.Nil(t, f)
			assert.ErrorIs(t, err, ErrNoValue)
		}()

		// Read goroutine 2
		wg.Add(1)
		go func() {
			defer wg.Done()

			f, err := cachedDS.getCache(dataQueryA)
			assert.Nil(t, f)
			assert.ErrorIs(t, err, ErrNoValue)
		}()

		wg.Wait()
	})

	t.Run("write to and read from cache concurrently", func(t *testing.T) {
		var wg sync.WaitGroup

		// Write goroutine 1
		wg.Add(1)
		go func() {
			defer wg.Done()

			f, err := cachedDS.saveCache(dataQueryA, mockFramer{frames: framesA}, nil)
			assert.NoError(t, err)
			assert.Equal(t, framesA, f.Frames())
		}()

		// Write goroutine 2
		wg.Add(1)
		go func() {
			defer wg.Done()

			f, err := cachedDS.saveCache(dataQueryB, mockFramer{frames: framesB}, nil)
			assert.NoError(t, err)
			assert.Equal(t, framesB, f.Frames())
		}()

		// Wait for writing goroutines
		wg.Wait()

		// Read goroutine 1
		wg.Add(1)
		go func() {
			defer wg.Done()

			f, err := cachedDS.getCache(dataQueryA)
			assert.NoError(t, err)
			assert.Equal(t, framesA, f.Frames())
		}()

		// Read goroutine 2
		wg.Add(1)
		go func() {
			defer wg.Done()

			f, err := cachedDS.getCache(dataQueryB)
			assert.NoError(t, err)
			assert.Equal(t, framesB, f.Frames())
		}()

		// Wait for reading goroutines
		wg.Wait()
	})

	t.Run("read from the cache concurrently", func(t *testing.T) {
		var wg sync.WaitGroup

		// Read goroutine 1
		wg.Add(1)
		go func() {
			defer wg.Done()

			f, err := cachedDS.getCache(dataQueryA)
			assert.NoError(t, err)
			assert.Equal(t, framesA, f.Frames())
		}()

		// Read goroutine 2
		wg.Add(1)
		go func() {
			defer wg.Done()

			f, err := cachedDS.getCache(dataQueryB)
			assert.NoError(t, err)
			assert.Equal(t, framesB, f.Frames())
		}()

		// Wait for reading goroutines
		wg.Wait()
	})
}
