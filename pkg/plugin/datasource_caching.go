package plugin

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/pkg/errors"
)

// CacheDuration is a constant that defines how long to keep cached elements before they are refreshed
const CacheDuration = time.Minute * 5

// CacheCleanupInterval is the interval at which the internal cache is cleaned / garbage collected
const CacheCleanupInterval = time.Minute * 10

// ErrNoValue is returned when a cached value is not available in the local cache
var ErrNoValue = errors.New("no cached value was found with that key")

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

// Make sure Datasource implements required interfaces.
var (
	_ backend.QueryDataHandler   = (*CachedDatasource)(nil)
	_ backend.CheckHealthHandler = (*CachedDatasource)(nil)
)

// The CachedDatasource wraps the Datasource type and stores an internal map, and responds to queries with cached data.
// If there is no cached data to respond with, the CachedDatasource forwards the request to the Datasource
type CachedDatasource struct {
	datasource Datasource
	cache      map[string]CachedResult
}

func (c *CachedDatasource) getCache(req backend.DataQuery) (dfutil.Framer, error) {
	key, err := getCacheKey(req)
	if err != nil {
		return nil, err
	}

	// Return cached value if it's there and it's not expired
	res, ok := c.cache[key]
	if ok {
		expired := res.ExpiresAt.Before(time.Now())
		if !expired {
			return res.Result, nil
		}
	}

	return nil, errors.Wrap(ErrNoValue, key)
}

func (c *CachedDatasource) saveCache(req backend.DataQuery, f dfutil.Framer, err error) (dfutil.Framer, error) {
	// don't store cached values if an error was returned
	if err != nil {
		return f, err
	}
	key, err := getCacheKey(req)
	if err != nil {
		return nil, err
	}

	c.cache[key] = newCachedResult(f)
	return f, err
}

// HandleRepositoriesQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleRepositoriesQuery(ctx context.Context, q *models.RepositoriesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleRepositoriesQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleIssuesQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleIssuesQuery(ctx context.Context, q *models.IssuesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleIssuesQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleCommitsQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleCommitsQuery(ctx context.Context, q *models.CommitsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleCommitsQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleTagsQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleTagsQuery(ctx context.Context, q *models.TagsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleTagsQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleReleasesQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleReleasesQuery(ctx context.Context, q *models.ReleasesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleReleasesQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleContributorsQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleContributorsQuery(ctx context.Context, q *models.ContributorsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleContributorsQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandlePullRequestsQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandlePullRequestsQuery(ctx context.Context, q *models.PullRequestsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandlePullRequestsQuery(ctx, q, req)
	return c.saveCache(req, f, err)

}

// HandleLabelsQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleLabelsQuery(ctx context.Context, q *models.LabelsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleLabelsQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandlePackagesQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandlePackagesQuery(ctx context.Context, q *models.PackagesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandlePackagesQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleMilestonesQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleMilestonesQuery(ctx context.Context, q *models.MilestonesQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleMilestonesQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleVulnerabilitiesQuery is the cache wrapper for the issue query handler
func (c *CachedDatasource) HandleVulnerabilitiesQuery(ctx context.Context, q *models.VulnerabilityQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleVulnerabilitiesQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleProjectsQuery is the cache wrapper for the project query handler
func (c *CachedDatasource) HandleProjectsQuery(ctx context.Context, q *models.ProjectsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleProjectsQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleStargazersQuery is the cache wrapper for the stargazer query handler
func (c *CachedDatasource) HandleStargazersQuery(ctx context.Context, q *models.StargazersQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleStargazersQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleWorkflowsQuery is the cache wrapper for the workflows query handler
func (c *CachedDatasource) HandleWorkflowsQuery(ctx context.Context, q *models.WorkflowsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleWorkflowsQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleWorkflowUsageQuery is the cache wrapper for the workflows usage query handler
func (c *CachedDatasource) HandleWorkflowUsageQuery(ctx context.Context, q *models.WorkflowUsageQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleWorkflowUsageQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// HandleWorkflowRunsQuery is the cache wrapper for the workflows runs query handler
func (c *CachedDatasource) HandleWorkflowRunsQuery(ctx context.Context, q *models.WorkflowRunsQuery, req backend.DataQuery) (dfutil.Framer, error) {
	if value, err := c.getCache(req); err == nil {
		return value, err
	}

	f, err := c.datasource.HandleWorkflowRunsQuery(ctx, q, req)
	return c.saveCache(req, f, err)
}

// CheckHealth forwards the request to the datasource and does not perform any caching
func (c *CachedDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return c.datasource.CheckHealth(ctx, req)
}

// QueryData forwards the request to the datasource and does not perform any caching
func (c *CachedDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return c.datasource.QueryData(ctx, req)
}

func getCacheKey(req backend.DataQuery) (string, error) {
	m := map[string]interface{}{
		"query":    req.JSON,
		"interval": req.Interval,
		"type":     req.QueryType,
		"time_range": backend.TimeRange{
			To:   req.TimeRange.To.Round(CacheDuration),
			From: req.TimeRange.From.Round(CacheDuration),
		},
	}
	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(m); err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := h.Write(b.Bytes()); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// Cleanup removes old cache keys
func (c *CachedDatasource) Cleanup() {
	for k, v := range c.cache {
		// If it's an expired value, then delete it
		if v.ExpiresAt.Before(time.Now()) {
			delete(c.cache, k)
		}
	}
}

func (c *CachedDatasource) startCleanup() {
	t := time.NewTicker(CacheCleanupInterval)

	for {
		<-t.C
		c.Cleanup()
	}
}

// WithCaching accepts a Client and returns a CachedClient which wraps the provided Client
func WithCaching(datasource Datasource) *CachedDatasource {
	c := &CachedDatasource{
		datasource: datasource,
		cache:      map[string]CachedResult{},
	}

	go c.startCleanup()

	return c
}
