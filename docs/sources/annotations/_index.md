---
title: 'Annotations'
menuTitle: Annotations
description: Using annotations for the GitHub data source plugin
aliases:
  - github
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
weight: 400
---

# Annotations

[Annotations](https://grafana.com/docs/grafana/latest/dashboards/annotations) let you extract data from a data source and use it to annotate a dashboard.

To create annotations, you need to specify at least the following two fields:

- A String field for the annotation text
- A Time field for the annotation time

Annotations overlay events on a graph.

![Annotations on a graph](/media/docs/grafana/data-sources/github/annotations.png)

With annotations, you can display the following GitHub resources on a graph:

- Commits
- Issues
- Pull requests
- Releases
- Tags


All annotations require that you select a field to display on the annotation, and a field that represents the time that the event occurred.

![Annotations editor](/media/docs/grafana/data-sources/github/annotations-editor.png)

If you want to add titles or tags to the annotations, you can add additional fields with the appropriate types.

For more information on how to configure a query, refer to [Built-in query editor](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/#built-in-query).
