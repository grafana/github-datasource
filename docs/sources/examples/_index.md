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

For documentation on importing dashboards, check out the documentation on [grafana.com](https://grafana.com/docs/grafana/latest/reference/export_import/#importing-a-dashboard).

The sample dashboard can be obtained from either of two places.

- Importing from the Grafana dashboards page [located here](https://grafana.com/grafana/dashboards/14000)

### Importing from the Grafana server WebUI

- From the GitHub data source configuration page:

  1. Navigate to "Connections" and then click "Data sources".
  2. Select the GitHub plugin and click the "Dashboards" tab.
  3. Click "Import".

  {{% admonition type="note" %}}
  If loading it from this repository, open Grafana and click "Import Dashboard". Copy the JSON in `./src/dashboards/dashboard.json`, and paste it into the "Import via panel json" box.
  {{% /admonition %}}

## 2. Play demo

You can also take a look at the live dashboards on Grafana Play.

{{< docs/play title="GitHub data source plugin demo" url="https://play.grafana.org/dashboards/f/bb613d16-7ee5-4cf4-89ac-60dd9405fdd7/demo-github" >}}
