# Grafana GitHub datasource

[![Known Vulnerabilities](https://snyk.io/test/github/grafana/github-datasource/badge.svg)](https://snyk.io/test/github/grafana/github-datasource)
[![Maintainability](https://api.codeclimate.com/v1/badges/30a924eb80d5f6b1cf9c/maintainability)](https://codeclimate.com/github/grafana/github-datasource/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/30a924eb80d5f6b1cf9c/test_coverage)](https://codeclimate.com/github/grafana/github-datasource/test_coverage)

The GitHub datasource allows GitHub API data to be visually represented in Grafana dashboards.

## Github API V4 (graphql)

This datasource uses the [`githubv4` package](https://github.com/shurcooL/githubv4), which is under active development.

## Features

### Backend

- [x] Releases
- [x] Commits
- [x] Repositories
- [x] Issues
- [x] Organizations
- [x] Labels
- [x] Milestones
- [x] Response Caching
- [ ] Deploys

### Frontend

- [x] Visualize queries
- [x] Template variables
- [x] Annotations

## Caching

Caching on this plugin is always enabled.

## Configuration

Options:

| Setting               | Required |
| --------------------- | -------- |
| Access token          | true     |
| Default Organization  | false    |
| Default Repository    | true     |
| Github Enterprise URL | false    |

To create a new Access Token, navigate to [Personal Access Tokens](https://github.com/settings/tokens) and create a click "Generate new token."

## Annotations

Annotations overlay events on a graph.

![Annotations on a graph](https://github.com/grafana/github-datasource/raw/main/docs/screenshots/annotations.png)

With annotations, you can display:

- Commits
- Issues
- Pull Requests
- Releases
- Tags

on a graph.

All annotations require that you select a field to display on the annotation, and a field that represents the time that the event occurred.

![Annotations editor](https://github.com/grafana/github-datasource/raw/main/docs/screenshots/annotations-editor.png)

## Variables

[Variables](https://grafana.com/docs/grafana/latest/variables/) allow you to substitute values in a panel with pre-defined values.

![Creating Variables](https://github.com/grafana/github-datasource/raw/main/docs/screenshots/variables-create.png)

You can reference them inside queries, allowing users to configure parameters such as `Query` or `Repository`.

![Using Variables inside queries](https://github.com/grafana/github-datasource/raw/main/docs/screenshots/using-variables.png)

## Macros

You can use the following macros in your queries

| Macro Name | Syntax                     | Description                                                          | Example                                                                              |
| ---------- | -------------------------- | -------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| multiVar   | `$__multiVar(prefix,$var)` | Expands a multi value variable into github query string              | `$__multiVar(label,$labels)` will expand into `label:first-label label:second-label` |
|            |                            | When using **all** in multi variable, use **\*** as custom all value |                                                                                      |

## Access Token Permissions

For all repositories:

- `public_repo`
- `repo:status`
- `repo_deployment`
- `read:packages`
- `user:read`
- `user:email`

An extra setting is required for private repositories

- `repo (Full control of private repositories)`

## Sample Dashboard

For documentation on importing dashboards, check out the documentation on [grafana.com](https://grafana.com/docs/grafana/latest/reference/export_import/#importing-a-dashboard)

The sample dashboard can be obtained from either of two places.

1. From the Grafana dashboards page [located here](https://grafana.com/grafana/dashboards/14000)

2. From this repository

If loading it from this repository, open Grafana and click "Import Dashboard".

Copy the JSON in `./src/dashboards/dashboard.json`, and paste it into the "Import via panel json" box.

## Frequently Asked Questions

- **I am using GitHub OAuth on Grafana. Can my users make requests with their individual GitHub accounts instead of a shared `access_token`?**

No. This requires changes in Grafana first. See [this issue](https://github.com/grafana/grafana/issues/26023) in the Grafana project.

- **Why does it sometimes take up to 5 minutes for my new pull request / new issue / new commit to show up?**

We have aggressive caching enabled due to GitHub's rate limiting policies. When selecting a time range like "Last hour", a combination of the queries for each panel and the time range is cached temporarily.

- **Why are there two selection options for Pull Requests and Issue times when creating annotations?**

There are two times that affect an annotation:

- The time range of the dashboard or panel
- The time that should be used to display the event on the graph

The first selection is used to filter the events that display on the graph. For example, if you select "closed at", only events that were "closed" in your dashboard's time range will be displayed on the graph.

The second selection is used to determine where on the graph the event should be displayed.

Typically these will be the same, however there are some cases where you may want them to be different.
