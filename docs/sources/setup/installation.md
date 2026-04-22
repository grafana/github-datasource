---
title: Install the GitHub data source plugin for Grafana
menuTitle: Install
description: Learn how to install the GitHub data source plugin for Grafana
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

# Install the GitHub data source plugin for Grafana

You can any of the following sets of steps to install the GitHub data source plugin for Grafana.

## Install from plugin catalog

1. Sign in to Grafana as a server administrator.
1. Click **Administration** > **Plugins and data** > **Plugins** in the side navigation menu to view all plugins.
1. Type **GitHub** in the Search box.
1. Click the **All** in the **State** filter option.
1. Click the plugin’s logo.
1. Click **Install**.

Find more information about the [plugin catalog](https://grafana.com/docs/grafana/latest/administration/plugin-management/#plugin-catalog).

## Install from grafana.com

Install the plugin from the [grafana.com plugins page](https://grafana.com/grafana/plugins/grafana-github-datasource/?tab=installation) using the instructions provided there. With this installation, you will get the latest published version of the plugin.

## Install from GitHub

1. Go to [Releases](https://github.com/grafana/github-datasource/releases/) on the GitHub project page.

1. Find the release you want to install.

1. Download the release by clicking the release asset called `grafana-github-datasource-<VERSION>.zip`. You may need to un-collapse the **Assets** section to see it.

1. Unarchive the plugin into the Grafana plugins directory:

   On Linux or macOS, run the following commands to extract the plugin:

   ```bash
   unzip grafana-github-datasource-<VERSION>.zip
   mv grafana-github-datasource /var/lib/grafana/plugins
   ```

   On Windows, run the following command to extract the plugin:

   ```powershell
   Expand-Archive -Path grafana-github-datasource-<VERSION>.zip -DestinationPath C:\grafana\data\plugins
   ```

1. Restart Grafana.

## Install using grafana-cli

If you are using `grafana-cli`, execute the following command to install the latest published version of the plugin:

1. Login to your machine as `root` user.

1. On Linux or macOS, open your terminal and run the following command:

   ```bash
   grafana-cli plugins install grafana-github-datasource
   ```

   On Windows, run the following command:

   ```shell
   grafana-cli.exe plugins install grafana-github-datasource
   ```

1. Then restart Grafana.

### Install custom version

If you need custom version of the plugin from GitHub, you can install it by running the following command:

```bash
grafana-cli --pluginUrl <ZIP_FILE_URL> plugins install grafana-github-datasource
```

For example, to install version 1.6.0 of the plugin on Linux or macOS:

```bash
grafana-cli --pluginUrl https://github.com/grafana/github-datasource/releases/download/v1.6.0/grafana-github-datasource-1.6.0.zip plugins install grafana-github-datasource
```

Or to install version 1.6.0 of the plugin on Windows:

```powershell
grafana-cli.exe --pluginUrl https://github.com/grafana/github-datasource/releases/download/v1.6.0/grafana-github-datasource-1.6.0.zip plugins install grafana-github-datasource
```
