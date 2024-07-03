---
slug: '/authentication'
title: 'Authentication'
menuTitle: Authentication
description: Authenticating the GitHub data source plugin
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

# Authentication

1. After completing the Access Token Permissions, navigate into Grafana UI and click on the Menu option on the top left.

1. Navigate to the Connections and then click on the Data Sources.

1. Click on the GitHub Data source plugin which you have installed.

1. Go to its settings tab and at the bottom, paste the GitHub token which you have created above.
   ![Configuring API Token](/media/docs/grafana/data-sources/github/github-plugin-confg-token.png)

1. Click “Save & Test” button and you should see a confirmation message similar to this, that indicates that the Plugin is successfully configured.

   ![Testing Data source](/media/docs/grafana/data-sources/github/github-plugin-config-success.png)

   {{% admonition type="note" %}}
   Incase, if there are some errors then please check the Grafana logs for troubleshooting.
   {{% /admonition %}}
