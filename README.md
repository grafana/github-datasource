# Grafana GitHub datasource

The GitHub datasource allows GitHub API data to be visually represented in Grafana dashboards.

## Features

### Backend
* [x] Releases
* [x] Commits
* [x] Repositories
* [x] Issues
* [x] Organizations
* [x] Labels
* [x] Milestones
* [x] Response Caching
* [ ] Deploys

### Frontend
* [x] Visualize queries
* [x] Template variables
* [x] Annotations

## Caching

Caching on this plugin is always enabled.

## Configuration

Options:

| Setting | Required |
|---------|----------|
| Access token | true |
| Default Organization | false |
| Default Repository | true |

To create a new Access Token, navigate to [Personal Access Tokens](https://github.com/settings/tokens) and create a click "Generate new token."

## Annotations

## Variables

## Access Token Permissions

For all repositories:
* `public_repo`
* `repo:status`
* `repo_deployment`
* `read:packages`

* `user:read`
* `user:email`

An extra setting is required for private repositories
* `repo (Full control of private repositories)`


## Frequently Asked Questions

* I am using GitHub OAuth on Grafana. Can my users make requests with their individual GitHub accounts instead of a shared `access_token`?

No. This requires changes in Grafana first. See [this issue](https://github.com/grafana/grafana/issues/26023) in the Grafana project.

* Why does it sometimes take up to 5 minutes for my new pull request / new issue / new commit to show up?

We have aggressive caching enabled due to GitHub's rate limiting policies. When selecting a time range like "Last hour", a combination of the queries for each panel and the time range is cached temporarily.

* I am trying to use a template variable in the "Query" field and it's not working

Template variables are currently not supported outside of the "Owner / Organization" and "Repository" fields.
