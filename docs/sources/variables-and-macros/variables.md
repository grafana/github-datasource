---
slug: '/variables'
title: 'Variables'
menuTitle: Variables
description: Using Variables for the GitHub data source plugin
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
weight: 201
---

# Variables

[Variable](https://grafana.com/docs/grafana/latest/variables/) is a placeholder for a value that can be used in queries for a dashboard. 

Variables allow you to create more interactive and dynamic dashboards by replacing hard-coded values with dynamic options. They are displayed as dropdown lists at the top of the dashboard, making it easy to change the data being displayed.

**Example**

Here is an example of creating a Dashboard variable:

![Creating Variables](/media/docs/grafana/data-sources/github/variables-create.png)

You can reference them inside queries, allowing users to configure parameters such as `Query` or `Repository`.

![Using Variables inside queries](/media/docs/grafana/data-sources/github/using-variables.png)