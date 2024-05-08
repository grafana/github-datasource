### Provisioning

[It’s possible to configure data sources using config files with Grafana’s provisioning system](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

#### With the [prom-operator](https://github.com/prometheus-operator/prometheus-operator)

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