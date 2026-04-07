---
title: GitHub data source plugin for Grafana
menuTitle: GitHub data source
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
weight: 10
---

# GitHub data source plugin for Grafana

The GitHub data source plugin for Grafana lets you to query the GitHub API in Grafana so you can visualize your GitHub repositories and projects.

Watch this video to learn more about setting up the Grafana GitHub data source plugin: {{< youtube id="DW693S3cO48">}}

{{< docs/play title="GitHub data source plugin demo" url="https://play.grafana.org/d/cdgx261sa1ypsa/3-single-repo-with-override-examples" >}}

## Query types

The plugin supports the following query types:

- Code Scan
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
- Workflow runs

## Supported features

With the plugin you can:

- Visualize queries
- Use template variables
- Configure Annotations
- Cache queries

## Caching

Caching on this plugin is always enabled.

{{< admonition type="note" >}}
To work around [GitHub's rate limiting](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28), the plugin caches requests aggressively.

This can mean that it takes up to five minutes for a new pull request, commit, or issue to show up in a query.
{{< /admonition >}}

## Requirements

To use the GitHub data source plugin, you will need:

- A free [GitHub](https://github.com/) or a [GitHub Enterprise](https://github.com/enterprise) account.
- Any of the following Grafana editions:
  - Grafana OSS server.
  - A [Grafana Cloud](https://grafana.com/pricing/) stack.
  - An on-premise Grafana Enterprise server with an [activated license](https://grafana.com/docs/grafana/latest/enterprise/license/activate-license/).

## Get started

- To start using the plugin, you need to generate an access token, then install and configure the plugin. To do this, refer to [Setup](./setup).
- To use variable and macros, for creating a dynamic dashboard, refer to [Variables and Macros](./variables-and-macros).
- To annotate the data by displaying the GitHub resources on the dashboard, refer to [Annotations](./annotations/).
- To quickly visualize GitHub data in Grafana, refer to [Sample dashboards](./sample-dashboards/).

## Get the most out of the plugin

- Add [Annotations](https://grafana.com/docs/grafana/latest/dashboards/annotations/)
- Configure and use [Templates and variables](https://grafana.com/docs/grafana/latest/variables/)
- Add [Transformations](https://grafana.com/docs/grafana/latest/panels/transformations/)

## Report issues

Use the [official GitHub repository](https://github.com/grafana/github-datasource/issues) to report issues, bugs, and feature requests.
