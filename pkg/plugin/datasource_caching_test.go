package plugin

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
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
