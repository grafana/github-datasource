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

1. Click on the GitHub data source plugin which you have installed.

1. Go to its settings tab and at the bottom, you will find the **Connection** section.

1. Paste the access token.
   ![Configuring API Token](/media/docs/grafana/data-sources/github/github-plugin-confg-token.png)

   (_Optional_): If you using the GitHub Enterprise, then select the **Enterprise** option inside the **Additional Settings** section and paste the URL of your GitHub Enterprise.

1. Click **Save & Test** button and you should see a confirmation dialog box that says "Data source is working".

{{< admonition type="tip" >}}
If you see errors, check the Grafana logs for troubleshooting.
{{< /admonition >}}
