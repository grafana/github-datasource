---
title: Create sample dashboards for the Grafana GitHub data source plugin for Grafana
menuTitle: Create sample dashboards
description: Learn how to import example dashboards into Grafana for use with the GitHub data source plugin for Grafana
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

# Create sample dashboards for the GitHub data source plugin for Grafana

This page explains how you can create a sample dashboard in Grafana to get started with the GitHub data source plugin. You can obtain these sample dashboards by:

- Using pre-configured dashboards
- Using play demo

## Use pre-configured dashboards

The pre-configured dashboards are ready-to-use and you only need to import them inside your Grafana server. There are two ways to use the pre-configure dashboards:

- Importing from the official Website
- Importing from the Grafana server WebUI

### Import from the Dashboards page on grafana.com

Import the [GitHub Default dashboard](https://grafana.com/grafana/dashboards/14000).

For instructions on how to import dashboards in Grafana, refer to [Import a dashboard](https://grafana.com/docs/grafana/latest/reference/export_import/#importing-a-dashboard).
The dashboard ID is `14000`.

### Import in the Grafana UI

To import a dashboard in the Grafana UI:

1. Go to **Connections** in the sidebar menu.
1. Under Connections, click **Data sources**.
1. Type `GitHub` in the search bar and select the GitHub data source.
1. Go to the **Dashboards** tab to view a list of pre-made dashboards.
1. Click **Import** to import the pre-made dashboard.

## Play demo

The Play demo dashboards provides a reference dashboard and allows you to create your own custom dashboards.

{{< docs/play title="GitHub data source plugin demo" url="https://play.grafana.org/dashboards/f/bb613d16-7ee5-4cf4-89ac-60dd9405fdd7/demo-github" >}}
