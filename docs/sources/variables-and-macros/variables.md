---
title: Create variable with GitHub data source plugin for Grafana
menuTitle: Variables
description: Learn about the variables you can use in the GitHub data source plugin for Grafana
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
weight: 301
---

# Create variable with GitHub data source plugin for Grafana

A [variable](https://grafana.com/docs/grafana/latest/variables/) is a placeholder for a value that you can use in dashboard queries.

Variables allow you to create more interactive and dynamic dashboards by replacing hard-coded values with dynamic options. They are displayed as dropdown lists at the top of the dashboard, making it easy to change the data being displayed.

**Example**

Here is an example of creating a dashboard variable:

{{< figure alt="Creating variables" src="/media/docs/grafana/data-sources/github/variables-create.png" >}}

You can reference them inside queries, allowing users to configure parameters such as `Query` or `Repository`.

{{< figure alt="Using variables inside queries" src="/media/docs/grafana/data-sources/github/using-variables.png" >}}
