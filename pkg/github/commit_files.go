package github

import (
	"context"
	"fmt"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/github-datasource/pkg/models"
)

// CommitFilesWrapper is a list of commit files returned by the GitHub API
type CommitFilesWrapper []*googlegithub.CommitFile

// Frames converts the list of commit files to a Grafana DataFrame
func (files CommitFilesWrapper) Frames() data.Frames {
	frame := data.NewFrame(
		"commit_files",
		data.NewField("path", nil, []string{}),
		data.NewField("additions", nil, []int64{}),
		data.NewField("deletions", nil, []int64{}),
		data.NewField("changes", nil, []int64{}),
		data.NewField("status", nil, []string{}),
		data.NewField("previous_filename", nil, []string{}),
	)

	for _, f := range files {
		frame.AppendRow(
			f.GetFilename(),
			int64(f.GetAdditions()),
			int64(f.GetDeletions()),
			int64(f.GetChanges()),
			f.GetStatus(),
			f.GetPreviousFilename(),
		)
	}

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	return data.Frames{frame}
}

// GetCommitFiles fetches the files changed in a specific commit.
// The GitHub REST API returns at most 300 files for a single commit.
func GetCommitFiles(ctx context.Context, client models.Client, opts models.CommitFilesOptions) (CommitFilesWrapper, error) {
	if opts.Owner == "" || opts.Repository == "" || opts.Ref == "" {
		return nil, nil
	}

	files, _, err := client.GetCommitFiles(ctx, opts.Owner, opts.Repository, opts.Ref, &googlegithub.ListOptions{
		PerPage: 300,
	})
	if err != nil {
		return nil, fmt.Errorf("getting commit files: owner=%s repo=%s sha=%s: %w", opts.Owner, opts.Repository, opts.Ref, err)
	}

	return CommitFilesWrapper(files), nil
}

// GetPullRequestFiles fetches all files changed in a pull request, handling pagination.
func GetPullRequestFiles(ctx context.Context, client models.Client, opts models.PullRequestFilesOptions) (CommitFilesWrapper, error) {
	if opts.Owner == "" || opts.Repository == "" || opts.PRNumber == 0 {
		return nil, nil
	}

	var allFiles []*googlegithub.CommitFile
	page := 1

	for {
		files, resp, err := client.ListPullRequestFiles(ctx, opts.Owner, opts.Repository, int(opts.PRNumber), &googlegithub.ListOptions{
			Page:    page,
			PerPage: 100,
		})
		if err != nil {
			return nil, fmt.Errorf("listing PR files: owner=%s repo=%s pr=%d page=%d: %w", opts.Owner, opts.Repository, opts.PRNumber, page, err)
		}

		allFiles = append(allFiles, files...)

		if resp == nil || resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}

	return CommitFilesWrapper(allFiles), nil
}

// CommitWithFiles holds a commit and the files changed in it
type CommitWithFiles struct {
	Commit Commit
	Files  []*googlegithub.CommitFile
}

// CommitsWithFiles is a list of commits each paired with their changed files
type CommitsWithFiles []CommitWithFiles

// Frames converts the list of commits-with-files to a flattened Grafana DataFrame.
// Each row represents one file change within a commit (one row per commit × file).
func (c CommitsWithFiles) Frames() data.Frames {
	frame := data.NewFrame(
		"commits",
		data.NewField("id", nil, []string{}),
		data.NewField("author", nil, []string{}),
		data.NewField("author_login", nil, []string{}),
		data.NewField("author_email", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("committed_at", nil, []time.Time{}),
		data.NewField("pushed_at", nil, []time.Time{}),
		data.NewField("message", nil, []string{}),
		data.NewField("file_path", nil, []string{}),
		data.NewField("file_additions", nil, []int64{}),
		data.NewField("file_deletions", nil, []int64{}),
		data.NewField("file_changes", nil, []int64{}),
		data.NewField("file_status", nil, []string{}),
		data.NewField("previous_filename", nil, []string{}),
	)

	for _, cwf := range c {
		for _, f := range cwf.Files {
			frame.AppendRow(
				cwf.Commit.OID,
				cwf.Commit.Author.Name,
				cwf.Commit.Author.User.Login,
				cwf.Commit.Author.Email,
				cwf.Commit.Author.User.Company,
				cwf.Commit.CommittedDate.Time,
				cwf.Commit.PushedDate.Time,
				string(cwf.Commit.Message),
				f.GetFilename(),
				int64(f.GetAdditions()),
				int64(f.GetDeletions()),
				int64(f.GetChanges()),
				f.GetStatus(),
				f.GetPreviousFilename(),
			)
		}
	}

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	return data.Frames{frame}
}

