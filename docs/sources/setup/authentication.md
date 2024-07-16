---
title: 'Configure authentication'
menuTitle: Configure authentication
description: Configure the GitHub data source plugin to authenticate to GitHub.
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

   {{< admonition type="note" >}}
   If you are using GitHub OAuth on Grafana, then it is not possible that the users to make requests with their individual GitHub accounts instead of a shared `access_token`. Please refer to [this issue](https://github.com/grafana/grafana/issues/26023) in the Grafana project.
   {{< /admonition >}}
