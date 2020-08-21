package github

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/dfutil"
	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// CacheDuration is a constant that defines how long to keep cached elements before they are refreshed
const CacheDuration = time.Minute

// CacheCleanupInterval is the interval at which the internal cache is cleaned / garbage collected
const CacheCleanupInterval = time.Minute * 10

// CachedResult is a value and a timestamp that defines when the cached value is no longer usable
type CachedResult struct {
	Result    dfutil.Framer
	ExpiresAt time.Time
}

func newCachedResult(f dfutil.Framer) CachedResult {
	return CachedResult{
		ExpiresAt: time.Now().Add(CacheDuration),
		Result:    f,
	}
}

// The CachedDatasource wraps the Datasource type and stores an internal map, and responds to queries with cached data.
// If there is no cached data to respond with, the CachedDatasource forwards the request to the Datasource
type CachedDatasource struct {
	datasource models.Datasource
	cache      map[string]CachedResult
}

func (c *CachedDatasource) cacheHandler(ctx context.Context, q *models.Query, req backend.DataQuery, handler models.QueryHandlerFunc) (dfutil.Framer, error) {
	key, err := getCacheKey(q, req)
	if err != nil {
		return nil, err
	}

	// Return cached value if it's there and it's not expired
	res, ok := c.cache[key]
	if ok && res.ExpiresAt.Before(time.Now()) {
		return res.Result, nil
	}

	// Otherwise run the handler

	framer, err := handler(ctx, q, req)
	if err != nil {
		return nil, err
	}

	c.cache[key] = newCachedResult(framer)

	return framer, nil
}

func (c *CachedDatasource) HandleIssuesQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return c.cacheHandler(ctx, q, req, c.datasource.HandleIssuesQuery)
}

func (c *CachedDatasource) HandleCommitsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return c.cacheHandler(ctx, q, req, c.datasource.HandleCommitsQuery)
}

func (c *CachedDatasource) HandleTagsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return c.cacheHandler(ctx, q, req, c.datasource.HandleTagsQuery)
}

func (c *CachedDatasource) HandleReleasesQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return c.cacheHandler(ctx, q, req, c.datasource.HandleReleasesQuery)
}

func (c *CachedDatasource) HandleContributorsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return c.cacheHandler(ctx, q, req, c.datasource.HandleContributorsQuery)
}

func (c *CachedDatasource) HandlePullRequestsQuery(ctx context.Context, q *models.Query, req backend.DataQuery) (dfutil.Framer, error) {
	return c.cacheHandler(ctx, q, req, c.datasource.HandlePullRequestsQuery)
}

func getCacheKey(query *models.Query, req backend.DataQuery) (string, error) {
	m := map[string]interface{}{
		"q":   query,
		"req": req,
	}
	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(m); err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(b.Bytes())

	return hex.EncodeToString(h.Sum(nil)), nil
}

// Cleanup removes old cache keys
func (c *CachedDatasource) Cleanup() {
	for k, v := range c.cache {
		if v.ExpiresAt.Before(time.Now()) {
			delete(c.cache, k)
		}
	}
}

func (c *CachedDatasource) startCleanup() {
	t := time.NewTicker(CacheCleanupInterval)

	for {
		select {
		case <-t.C:
			c.Cleanup()
		}
	}
}

// WithCaching accepts a Client and returns a CachedClient which wraps the provided Client
func WithCaching(datasource models.Datasource) *CachedDatasource {
	c := &CachedDatasource{
		datasource: datasource,
		cache:      map[string]CachedResult{},
	}

	go c.startCleanup()

	return c
}
