---
title: 'Installation'
menuTitle: Installation
description: Installation of the GitHub data source plugin
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
weight: 101
---

# Installing GitHub data source plugin

There are multiple ways to install this plugin into your Grafana instance

## Install from plugin catalog

Install the plugin directly within the Grafana web UI using the [plugin catalog](https://grafana.com/docs/grafana/latest/administration/plugin-management/#plugin-catalog).

## Install from grafana.com

Install the plugin from the [grafana.com plugins page](https://grafana.com/grafana/plugins/grafana-github-datasource/?tab=installation) using the instructions provided there. With this installation, you will get the latest published version of the plugin.

## Install from GitHub

Download the required version of the release zip file from [github](https://github.com/grafana/github-datasource/releases/) and extract it into your grafana plugin folder. Then restart Grafana.

## Install using grafana-cli

If you are using grafana-cli, execute the following command to install the latest published version of the plugin

```bash
grafana-cli plugins install grafana-github-datasource
```

If you need custom version of the plugin from github, you can install using the following command.

```bash
grafana-cli --pluginUrl <ZIP_FILE_URL> plugins install grafana-github-datasource
```

Example:

```bash
grafana-cli --pluginUrl https://github.com/grafana/github-datasource/releases/download/v1.6.0/grafana-github-datasource-1.6.0.zip plugins install grafana-github-datasource
```
