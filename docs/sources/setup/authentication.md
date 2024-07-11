---
title: 'Configure authentication'
menuTitle: Configure authentication
description: Configure the GitHub data source plugin to authenticate to GitHub.
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
weight: 103
---

# Configure authentication

1. After completing the **Access Token Permissions** in GitHub, navigate into Grafana and click on the menu option on the top left.

1. Browse to the **Connections** menu and then click on the **Data sources**.

1. Click on the GitHub data source plugin which you have installed.

1. Go to its settings tab and at the bottom, paste the GitHub token which you have created above.
   ![Configuring API Token](/media/docs/grafana/data-sources/github/github-plugin-confg-token.png)

1. Click **Save & Test** button and you should see a confirmation dialog box that says "Data source is working".

   {{< admonition type="note" >}}
   If you see errors, check the Grafana logs for troubleshooting.
   {{< /admonition >}}
