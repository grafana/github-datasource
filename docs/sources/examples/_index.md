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

Here are some examples dashboards that you can use for your use case.

## Pre-configure dashboards

For documentation on importing dashboards, check out the documentation on [grafana.com](https://grafana.com/docs/grafana/latest/reference/export_import/#importing-a-dashboard).

The sample dashboard can be obtained from either of two places.

- Importing from the Grafana dashboards page [located here](https://grafana.com/grafana/dashboards/14000)

OR

- From the GitHub data source configuration page:

  1. Navigate to "Connections" and then click "Data sources".
  1. Select the GitHub plugin and click the "Dashboards" tab.
  1. Click "Import".

  {{% admonition type="note" %}}
  If loading it from this repository, open Grafana and click "Import Dashboard". Copy the JSON in `./src/dashboards/dashboard.json`, and paste it into the "Import via panel json" box.
  {{% /admonition %}}

## Play demo

You can also take a look at the live dashboards on play using this [link](https://play.grafana.org/dashboards/f/bb613d16-7ee5-4cf4-89ac-60dd9405fdd7/demo-github).
