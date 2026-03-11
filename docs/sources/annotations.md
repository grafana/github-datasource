---
aliases:
  - ./annotations/
title: GitHub annotations
menuTitle: Annotations
description: Use annotations to overlay GitHub events on Grafana dashboard panels
keywords:
  - data source
  - github
  - annotations
  - commits
  - issues
  - pull requests
  - releases
  - tags
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 400
review_date: "2026-03-11"
---

# GitHub annotations

[Annotations](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/) let you extract data from the GitHub data source and overlay events on dashboard panels.

## Required fields

To create annotations, you need to specify at least the following two fields:

- A **String** field for the annotation text.
- A **Time** field for the annotation time.

If you want to add titles or tags to the annotations, you can add additional fields with the appropriate types.

## Supported resources

You can create annotations from the following GitHub resources:

- Commits
- Issues
- Pull requests
- Releases
- Tags

## Configure an annotation query

To configure an annotation query:

1. Navigate to **Dashboard settings** > **Annotations**.
1. Click **Add annotation query**.
1. Select the GitHub data source.
1. Select the query type and configure the query options.
1. Select the field to display on the annotation and the field that represents the time the event occurred.

For more information on annotation configuration, refer to [Built-in query editor](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/#built-in-query).

## Pull request and issue annotation times

When you create annotations for pull requests and issues, there are two time selections. This is because two different times affect an annotation:

- **Filter time:** The time range of the dashboard or panel, used to filter which events display on the graph.
- **Display time:** The time that determines where on the graph the event appears.

The first selection filters the events that display on the graph. For example, if you select `ClosedAt`, only events that were closed within your dashboard's time range are displayed.

The second selection determines where on the graph the event is positioned.

Typically these are the same, however there are cases where you may want them to be different. For example, you might want to filter by `ClosedAt` but display the annotation at the `CreatedAt` time.
