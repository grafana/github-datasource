---
title: 'Examples dashboards'
menuTitle: Examples dashboards
description: Examples dashboards
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
weight: 500
---

# Example dashboard

There are multiple ways to use some of the examples dashboards that you can use for your use case.

1. Using pre-configured dashboards
1. Using play demo

## 1. Using pre-configure dashboards

There are two ways to use the pre-configure dashboards:

- Importing from the official Website
- Importing from the Grafana server WebUI

### Importing from the official Website

Import the dashboards from the official page [located here](https://grafana.com/grafana/dashboards/14000).

For instructions on how to import dashboards in Grafana, see [Import a dashboard](https://grafana.com/docs/grafana/latest/reference/export_import/#importing-a-dashboard) and replace the ID with `1400`.

### Importing from the Grafana server WebUI

To view a list of pre-made GitHub dashboards do the following:

1. Go to **Connections** in the sidebar menu.
1. Under Connections, click **Data sources**.
1. Type `GitHub` in the search bar and select the GitHub data source.
1. Go to the **Dashboards** tab to view a list of pre-made dashboards.
1. Click **Import** to import the pre-made dashboard.

## 2. Play demo

You can also take a look at the live dashboards on Grafana Play.

{{< docs/play title="GitHub data source plugin demo" url="https://play.grafana.org/dashboards/f/bb613d16-7ee5-4cf4-89ac-60dd9405fdd7/demo-github" >}}
