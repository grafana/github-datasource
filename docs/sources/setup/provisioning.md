---
title: 'Provisioning'
menuTitle: Provisioning
description: Provisioning the GitHub data source plugin
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
weight: 104
---

# Provisioning

You can define and configure the GitHub data source in YAML files with Grafana provisioning.
For more information about provisioning a data source, and for available configuration options, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

**Example**

Here is an example to provision the Data source while using the [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator)

```yaml
promop:
  grafana:
    additionalDataSources:
      - name: GitHub Repo Insights
        type: grafana-github-datasource
        jsonData:
          owner: ''
          repository: ''
        secureJsonData:
          accessToken: '<github api token>'
```
