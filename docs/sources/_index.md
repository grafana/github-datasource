---
aliases:
  - ./sample-dashboards/
title: GitHub data source plugin for Grafana
menuTitle: GitHub data source
description: The GitHub data source plugin lets you visualize GitHub data in Grafana dashboards.
keywords:
  - data source
  - github
  - github repository
  - API
  - pull requests
  - issues
  - workflows
  - commits
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 10
review_date: "2026-03-11"
---

# GitHub data source plugin for Grafana

The GitHub data source plugin lets you query the GitHub API so you can visualize and monitor your GitHub repositories, organizations, and projects in Grafana dashboards.

## Requirements

To use the GitHub data source plugin, you need:

- A [GitHub](https://github.com/) account or a [GitHub Enterprise](https://github.com/enterprise) account.
- Grafana version 10.4.8 or later (OSS, Enterprise, or Cloud).

## Key capabilities

The GitHub data source supports:

- **Query 19 resource types:** Commits, issues, pull requests, pull request reviews, workflows, workflow runs, deployments, code scanning alerts, and more.
- **Template variables and macros:** Create dynamic, reusable dashboards with variable-driven queries and built-in macros.
- **Annotations:** Overlay GitHub events (commits, issues, pull requests, releases, tags) on dashboard panels.
- **Alerting:** Create alert rules based on GitHub query results.
- **Built-in caching:** Automatic request caching to handle GitHub API rate limits.

## Get started

The following pages help you set up and use the GitHub data source:

- [Configure the GitHub data source](https://grafana.com/docs/plugins/grafana-github-datasource/latest/configure/)
- [GitHub query editor](https://grafana.com/docs/plugins/grafana-github-datasource/latest/query-editor/)
- [Template variables](https://grafana.com/docs/plugins/grafana-github-datasource/latest/template-variables/)
- [Annotations](https://grafana.com/docs/plugins/grafana-github-datasource/latest/annotations/)
- [Troubleshooting](https://grafana.com/docs/plugins/grafana-github-datasource/latest/troubleshooting/)

## Pre-built dashboards

The plugin includes a pre-built dashboard that you can import to get started quickly.

### Import from grafana.com

Import the [GitHub Default dashboard](https://grafana.com/grafana/dashboards/14000) (dashboard ID `14000`).

For instructions on importing dashboards, refer to [Import a dashboard](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/import-dashboards/).

### Import in the Grafana UI

To import the dashboard in the Grafana UI:

1. Click **Connections** in the left-side menu.
1. Click **Data sources**.
1. Select your GitHub data source.
1. Click the **Dashboards** tab.
1. Click **Import** next to the dashboard you want to use.

### Play demo

{{< docs/play title="GitHub data source plugin demo" url="https://play.grafana.org/dashboards/f/bb613d16-7ee5-4cf4-89ac-60dd9405fdd7/demo-github" >}}

## Caching

Caching is always enabled for this plugin.

{{< admonition type="note" >}}
To stay within [GitHub's rate limits](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api), the plugin caches all API responses. Cached data may be up to five minutes old, so recent changes to pull requests, commits, or issues may not appear immediately.
{{< /admonition >}}

## Plugin updates

Always ensure that your plugin version is up-to-date so you have access to all current features and improvements. Navigate to **Plugins and data** > **Plugins** to check for updates.

{{< admonition type="note" >}}
Plugins are automatically updated in Grafana Cloud.
{{< /admonition >}}

## Related resources

- [GitHub REST API documentation](https://docs.github.com/en/rest)
- [Grafana community forum](https://community.grafana.com/)
- [Report issues on GitHub](https://github.com/grafana/github-datasource/issues)
