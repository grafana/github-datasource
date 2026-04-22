---
title: Provisioning the GitHub data source in Grafana
menuTitle: Provisioning
description: Provisioning the GitHub data source plugin
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
weight: 104
---

# Provisioning the GitHub data source in Grafana

You can define and configure the GitHub data source in YAML files with Grafana provisioning. For more information about provisioning a data source, and for available configuration options, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

**Example**

```yaml
apiVersion: 1

datasources:
  - name: GitHub (Personal Access Token)
    type: grafana-github-datasource
    jsonData:
      selectedAuthType: personal-access-token
    secureJsonData:
      accessToken: <your_access_token>

  - name: GitHub (App)
    type: grafana-github-datasource
    jsonData:
      selectedAuthType: github-app
      appId: <your_app_id>
      installationId: <your_installation_id>
    secureJsonData:
      privateKey: <your_private_key>
```
