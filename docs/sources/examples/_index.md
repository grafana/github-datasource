---
title: Create a sample dashboards for the Grafana GitHub data source plugin
menuTitle: Sample dashboards
description: Learn how to import pre-made example dashboards
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

# Create a sample dashboards the GitHub data source plugin using Grafana

This page explains how you can create a sample dashboard in Grafana to get started with the GitHub data source plugin. You can obtain these sample dashboards by:

- Using pre-configured dashboards
- Using play demo

## Using pre-configure dashboards

The pre-configured dashboards are ready-to-use and you only need to import them inside your Grafana server. There are two ways to use the pre-configure dashboards:

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

## Play demo

The Play demo dashboards prvoide a reference, allows you to create your own custom dashboards.

{{< docs/play title="GitHub data source plugin demo" url="https://play.grafana.org/dashboards/f/bb613d16-7ee5-4cf4-89ac-60dd9405fdd7/demo-github" >}}
