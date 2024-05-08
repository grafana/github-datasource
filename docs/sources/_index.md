title: GitHub data source plugin for Grafana
menuTitle: GitHub data source
description: This document introduces the GitHub data source
aliases:
  - github 
keywords:
  - data source
  - github
  - github repository
  - api
labels:
  products:
    - oss
    - enterprise
    - grafana cloud
weight: 10
---

# Overview

The GitHub datasource allows GitHub API data to be visually represented in Grafana dashboards.

## GitHub API V4 (GraphQL)

This datasource uses the [`githubv4` package](https://github.com/shurcooL/githubv4), which is under active development.

## Key Features

### Backend

- [x] Releases
- [x] Commits
- [x] Repositories
- [x] Stargazers
- [x] Issues
- [x] Organizations
- [x] Labels
- [x] Milestones
- [x] Response Caching
- [x] Projects
- [x] Workflows
- [ ] Deploys

### Frontend

- [x] Visualize queries
- [x] Template variables
- [x] Annotations

### Caching

Caching on this plugin is always enabled.


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