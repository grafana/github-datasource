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

1. Login to Grafana Web UI as *Server Admin*
1. Click **Administration** > **Plugins and data** > **Plugins** in the side navigation menu to view all plugins.
1. Type **GitHub** in the Search box
1. Click the **All** in the **State** filter option.
1. Click the plugin’s logo.
1. Click Install.

Find more information about the [plugin catalog](https://grafana.com/docs/grafana/latest/administration/plugin-management/#plugin-catalog).

## Install from grafana.com

Install the plugin from the [grafana.com plugins page](https://grafana.com/grafana/plugins/grafana-github-datasource/?tab=installation) using the instructions provided there. With this installation, you will get the latest published version of the plugin.

## Install from GitHub

1. Go to [Releases](https://github.com/grafana/github-datasource/releases/) on the GitHub project page.

1. Find the release you want to install

1. Download the release by clicking the release asset called `grafana-github-datasource-<version>.zip`. You may need to un-collapse the **Assets** section to see it.

1. Unarchive the plugin into the Grafana plugins directory:

   In Linux/macOS, you can use the following command to extract the plugin:

   ```bash
   unzip grafana-github-datasource-<version>.zip
   mv grafana-github-datasource /var/lib/grafana/plugins
   ```

   In Windows, you can use the following command to extract the plugin:

   ```shell
   Expand-Archive -Path grafana-github-datasource-<version>.zip -DestinationPath C:\grafana\data\plugins
   ```

1. Then restart Grafana.

## Install using grafana-cli

If you are using grafana-cli, execute the following command to install the latest published version of the plugin

1. Login to your machine as root user.

1. In Linux/macOS, ypen your terminal and write the command:

   ```bash
   grafana-cli plugins install grafana-github-datasource
   ```

   Similarly, if you are using Windows machine, use the following command:

   ```shell
   grafana-cli.exe plugins install grafana-github-datasource
   ```

1. Then restart Grafana.

### Install custom version

If you need custom version of the plugin from github, you can install using the following command:

```bash
grafana-cli --pluginUrl <ZIP_FILE_URL> plugins install grafana-github-datasource
```

Example:

1. In Linux/macOS:

   ```bash
   grafana-cli --pluginUrl https://github.com/grafana/github-datasource/releases/download/v1.6.0/grafana-github-datasource-1.6.0.zip plugins install grafana-github-datasource
   ```

1. In Windows:

   ```shell
   grafana-cli.exe --pluginUrl https://github.com/grafana/github-datasource/releases/download/v1.6.0/grafana-github-datasource-1.6.0.zip plugins install grafana-github-datasource
   ```

1. Then restart Grafana.
