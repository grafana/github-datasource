---
description: The GitHub data source lets you visualize GitHub API data in Grafana dashboards.
keywords:
  - data source
  - github
  - github repository
  - API
labels:
  products:
    - oss
    - enterprise
    - cloud
menuTitle: GitHub data source
title: GitHub data source plugin for Grafana
weight: 10
---

# Overview

The GitHub datasource allows GitHub API data to be visually represented in Grafana dashboards.

{{< docs/play title="GitHub Data source plugin demo" url="https://play.grafana.org/d/cdgx261sa1ypsa/3-single-repo-with-override-examples" >}}

## GitHub API V4 (GraphQL)

This datasource uses the [`githubv4` package](https://github.com/shurcooL/githubv4), which is under active development.

## Key Features

The Grafana GitHub Data source plugin provides several features, allow you to track and analyze GitHub data directly within Grafana, providing insights and visualizations for your GitHub repositories and projects.

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

## Requirements

The GitHub data source plugin has the following requirements:

- A GitHub account.
- Any free or paid Grafana Cloud plan or an activated on-prem Grafana Enterprise license. Contracted Cloud customers should refer to their agreement.
- Port 3000 enabled.
