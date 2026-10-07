package github

import (
	"context"
	"fmt"
	"strings"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type codeSearchClient interface {
	SearchCode(context.Context, string, *googlegithub.SearchOptions) (*googlegithub.CodeSearchResult, *googlegithub.Response, error)
}

// CodeSearch is a paginated GitHub code-search result.
type CodeSearch struct {
	Items    []*googlegithub.CodeResult
	Total    int
	Complete bool
}

// Frames converts code-search results into a table-friendly frame.
func (search CodeSearch) Frames() data.Frames {
	frame := data.NewFrame(
		"code_search",
		data.NewField("file_name", nil, []*string{}),
		data.NewField("path", nil, []*string{}),
		data.NewField("sha", nil, []*string{}),
		data.NewField("file_url", nil, []*string{}),
		data.NewField("repository_name", nil, []*string{}),
		data.NewField("full_name", nil, []*string{}),
		data.NewField("repository_url", nil, []*string{}),
		data.NewField("fragments", nil, []string{}),
		data.NewField("match_count", nil, []int64{}),
		data.NewField("total_count", nil, []int64{}),
		data.NewField("complete", nil, []bool{}),
	)
	for _, item := range search.Items {
		var repositoryName, repositoryFullName, repositoryURL *string
		if item.Repository != nil {
			repositoryName = item.Repository.Name
			repositoryFullName = item.Repository.FullName
			repositoryURL = item.Repository.HTMLURL
		}
		fragments := make([]string, 0, len(item.TextMatches))
		for _, match := range item.TextMatches {
			if match.Fragment != nil {
				fragments = append(fragments, match.GetFragment())
			}
		}
		frame.AppendRow(
			item.Name, item.Path, item.SHA, item.HTMLURL,
			repositoryName, repositoryFullName, repositoryURL,
			strings.Join(fragments, "\n"), int64(len(item.TextMatches)), int64(search.Total), search.Complete,
		)
	}
	return data.Frames{frame}
}

func getCodeSearch(ctx context.Context, client codeSearchClient, query, exactPath string, includeTextMatches, requireComplete bool) (CodeSearch, error) {
	if strings.TrimSpace(query) == "" {
		return CodeSearch{}, fmt.Errorf("code search query is required")
	}
	search := CodeSearch{Complete: true}
	rawResultCount := 0
	for page := 1; page != 0; {
		result, response, err := client.SearchCode(ctx, query, &googlegithub.SearchOptions{
			TextMatch:   includeTextMatches,
			ListOptions: googlegithub.ListOptions{Page: page, PerPage: 100},
		})
		if err != nil {
			return CodeSearch{}, err
		}
		if result == nil {
			return CodeSearch{}, fmt.Errorf("code search returned no result")
		}
		search.Total = result.GetTotal()
		search.Complete = search.Complete && !result.GetIncompleteResults()
		rawResultCount += len(result.CodeResults)
		for _, item := range result.CodeResults {
			if exactPath == "" || item.GetPath() == exactPath {
				search.Items = append(search.Items, item)
			}
		}
		page = response.NextPage
	}
	if rawResultCount < search.Total {
		search.Complete = false
	}
	if requireComplete && !search.Complete {
		return CodeSearch{}, fmt.Errorf("code search returned incomplete results: received %d of %d", rawResultCount, search.Total)
	}
	return search, nil
}
