---
description: The GitHub data source lets you visualize GitHub data in Grafana dashboards.
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

# GitHub data source plugin for Grafana

The Grafana GitHub data source plugin provides, allow you to track and analyze GitHub data directly within Grafana, providing insights and visualizations for your GitHub repositories and projects.

{{< docs/play title="GitHub data source plugin demo" url="https://play.grafana.org/d/cdgx261sa1ypsa/3-single-repo-with-override-examples" >}}

The plugin provide the features listed below:

## Query types

- Commits
- Issues
- Contributors
- Tags
- Releases
- Pull requests
- Labels
- Repositories
- Milestones
- Packages
- Vulnerabilities
- Projects
- Stargazers
- Workflows
- Workflow usage

## Supported features

- Visualize queries
- Template variables
- Annotations
- Query caching

## Caching

Caching on this plugin is always enabled.

{{< admonition type="note" >}}
Sometimes it may take up to 5 minutes as we have aggressive caching enabled due to [GitHub's rate limiting](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28) policies. When selecting a time range like "Last hour", a combination of the queries for each panel and the time range is cached temporarily.
{{< /admonition >}}

## Requirements

The GitHub data source plugin has the following requirements:

- A free [GitHub](https://github.com/) or a [GitHub Enterprise](https://github.com/enterprise) account.
- Any of the following Grafana flavours:
  - Grafana OSS
  - Free or paid [Grafana Cloud](https://grafana.com/pricing/) server
  - An [activated on-prem Grafana Enterprise](https://grafana.com/docs/grafana/latest/enterprise/license/activate-license/) server.

## Get the most out of the plugin

- Add [Annotations](https://grafana.com/docs/grafana/latest/dashboards/annotations/)
- Configure and use [Templates and variables](https://grafana.com/docs/grafana/latest/variables/)
- Add [Transformations](https://grafana.com/docs/grafana/latest/panels/transformations/)

## Reporting issues

Use the [official GitHub repository](https://github.com/grafana/github-datasource/issues) to report issues, bugs and feature requests.

### Further reading

- [GitHub v4 client library](https://github.com/shurcooL/githubv4)
- [Manage personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
