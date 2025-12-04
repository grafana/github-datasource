---
title: Configure the GitHub data source plugin for Grafana
menuTitle: Configure
description: Configure the GitHub data source plugin to authenticate to GitHub
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
weight: 103
---

# Configure the GitHub data source plugin for Grafana

1. After creating the **access token** in GitHub, navigate into Grafana and click on the menu option on the top left.

1. Browse to the **Connections** menu and then click on the **Data sources**.

1. Click on the **Add new data source** button

1. Click on the GitHub data source plugin which you have installed.

1. Go to its settings tab and at the bottom, you will find the **Authentication** section.

1. Paste the access token.
   {{< figure alt="Configuring API Token" src="/media/docs/grafana/data-sources/github/github-plugin-confg-token.png" >}}

1. (_Optional_): If you using the GitHub Enterprise Server, then select the **Enterprise Server** option inside the **Connection** section and paste the URL of your GitHub Enterprise Server.

1. Click **Save & Test** button and you should see a confirmation dialog box that says "Data source is working".

**Private data source connect** - _Only for Grafana Cloud users._

| Setting                         | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Private data source connect** | Establishes a private, secured connection between a Grafana Cloud stack and data sources within a private network. Use the drop-down to locate the PDC URL. For setup instructions, refer to [Private data source connect (PDC)](/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/) and [Configure PDC](/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/configure-pdc/#configure-grafana-private-data-source-connect-pdc). Click **Manage private data source connect** to open your PDC connection page and view your configuration details. |


{{< admonition type="tip" >}}
If you see errors, check the Grafana logs for troubleshooting.
{{< /admonition >}}
