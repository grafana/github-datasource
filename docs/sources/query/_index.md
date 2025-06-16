---
title: 'Query'
menuTitle: Query
description: Query
hero:
  title: Github data source plugin for Grafana
  level: 1
  width: 110
  image: https://raw.githubusercontent.com/grafana/github-datasource/refs/heads/main/src/img/github.svg
  height: 110
  description: The Github data source plugin allows you to query and visualize data from Github.
keywords:
  - github
  - ci
  - cd
query_types:
  title_class: pt-0 lh-1
  items:
    - title: Repositories
      description: List the github repositories for an organization
      href: /docs/grafana-github-datasource/latest/query/repositories/
    - title: Issues
      description: List the issues created for a github repo
      href: /docs/grafana-github-datasource/latest/query/issues/
    - title: Pull Requests
      description: List the pull requests created for a github repo
      href: /docs/grafana-github-datasource/latest/query/pull-requests/
labels:
  products:
    - oss
weight: 300
---

<!-- markdownlint-configure-file { "MD013": false, "MD033": false, "MD025": false, "MD034": false } -->

{{< docs/hero-simple key="hero" >}}

<hr style="margin-bottom:30px"/>

## 🎯 Supported query types

{{< card-grid key="query_types" type="simple" >}}

### Other query types

The plugin also supports the following query types:

- Commits
- Contributors
- Tags
- Releases
- Labels
- Milestones
- Packages
- Vulnerabilities
- Projects
- Stargazers
- Workflows
- Workflow usage
- Workflow runs
