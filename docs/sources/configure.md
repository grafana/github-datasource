---
aliases:
  - ./setup/
  - ./setup/installation/
  - ./setup/token/
  - ./setup/datasource/
  - ./setup/provisioning/
title: Configure the GitHub data source
menuTitle: Configure
description: Configure the GitHub data source to connect Grafana to GitHub
keywords:
  - data source
  - github
  - configuration
  - authentication
  - personal access token
  - github app
  - provisioning
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 100
review_date: "2026-03-11"
---

# Configure the GitHub data source

This document explains how to configure the GitHub data source for Grafana.

## Before you begin

To configure the data source, you need:

- **Grafana permissions:** Organization administrator role.
- **GitHub credentials** (one of the following):
  - A personal access token (classic or fine-grained) with the [required scopes](#personal-access-token-permissions).
  - A registered GitHub App with the **App ID**, **Installation ID**, and **private key**. Refer to [Register and configure a GitHub App](#register-and-configure-a-github-app) for setup steps.
- **GitHub Enterprise Server URL** (if applicable): The URL of your GitHub Enterprise Server instance.

### Security best practices

When creating credentials for the GitHub data source, follow the principle of least privilege:

- **Prefer GitHub Apps or fine-grained tokens over classic PATs.** GitHub Apps and fine-grained tokens let you restrict access to specific repositories and grant only the permissions the plugin needs. Classic personal access tokens grant access to all repositories the user can access.
- **Grant only read-only permissions.** The data source only reads data from GitHub. Never grant write access unless you have a specific requirement.
- **Scope tokens to the repositories you need.** Avoid granting organization-wide access when you only query a few repositories.
- **Rotate credentials regularly.** Set an expiration on personal access tokens and rotate them on a regular schedule. GitHub Apps use short-lived tokens that rotate automatically.

Grafana stores all credentials (access tokens, private keys) as encrypted secure JSON data. Credentials are never exposed in API responses or log output.

## Add the data source

To add the GitHub data source:

1. Click **Connections** in the left-side menu.
1. Click **Add new connection**.
1. Type `GitHub` in the search bar.
1. Select **GitHub**.
1. Click **Add new data source**.

## Configure settings

The following table describes the GitHub-specific configuration settings.

| Setting | Description |
|---------|-------------|
| **GitHub License Type** | Select your GitHub plan: **Free, Pro & Team**, **Enterprise Cloud**, or **Enterprise Server**. |
| **GitHub Enterprise Server URL** | The URL of your GitHub Enterprise Server instance. Only visible when **Enterprise Server** is selected as the GitHub license. |

### Private data source connect

{{< admonition type="note" >}}
Private data source connect is available for Grafana Cloud users only.
{{< /admonition >}}

Private data source connect (PDC) establishes a private, secured connection between a Grafana Cloud stack and data sources within a private network. Use the drop-down to select a PDC connection.

Click **Manage private data source connect** to go to your PDC connection page, where you can find your PDC configuration details.

For setup instructions, refer to [Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/).

## Authentication

The GitHub data source supports two authentication methods: personal access tokens and GitHub Apps.

### Personal access token

You can authenticate with either a classic personal access token or a fine-grained personal access token.

#### Create a classic personal access token

1. Sign in to your GitHub account.
1. Navigate to [Personal access tokens](https://github.com/settings/tokens) and click **Generate new token**.
1. Select **personal access token (classic)**.
1. Assign the [required permissions](#personal-access-token-permissions).
1. Click **Generate Token**.
1. Copy the token and paste it into the **Personal Access Token** field in the data source settings.

For more information, refer to the [GitHub personal access token documentation](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).

#### Create a fine-grained personal access token

1. Sign in to your GitHub account.
1. Navigate to [Fine-grained personal access tokens](https://github.com/settings/tokens?type=beta) and click **Generate new token**.
1. Provide a name for the token.
1. Assign the required repository access and `read-only` permissions.
1. Click **Generate token**.
1. Copy the token and paste it into the **Personal Access Token** field in the data source settings.

For more information, refer to the [GitHub fine-grained personal access token documentation](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token).

#### Personal access token permissions

The following scopes are required for classic personal access tokens:

| Scope | Purpose |
|-------|---------|
| `public_repo` | Access public repositories. |
| `repo:status` | Access commit status. |
| `repo_deployment` | Access deployment status. |
| `read:packages` | Read packages. |
| `read:user` | Read user profile data. |
| `user:email` | Access user email addresses. |
| `read:org` | Read organization membership. |
| `read:project` | Read project data. |
| `repo` | Full control of private repositories. Required only if you need to query private repositories. |

### GitHub App

GitHub App authentication provides better security and fine-grained access to resources compared to personal access tokens.

#### Register and configure a GitHub App

1. Register a new GitHub App by following the [GitHub App documentation](https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app).
1. After registering the app, generate a private key for authentication.
1. Note the **App ID** assigned to your GitHub App.
1. [Install the GitHub App](https://docs.github.com/en/apps/using-github-apps/installing-your-own-github-app) on your GitHub account or organization.
1. Note the **Installation ID** after completing the installation.
1. In the Grafana data source settings, provide the **App ID**, **Installation ID**, and **Private Key** in the appropriate fields.

{{< admonition type="note" >}}
To find your installation ID, navigate to **Settings** > **Installed GitHub Apps** > **Configure**. The installation ID is the number at the end of the URL: `https://github.com/settings/installations/<INSTALLATION_ID>`.
{{< /admonition >}}

#### GitHub App permissions

The following repository permissions are required:

| Permission | Access level |
|------------|-------------|
| **Metadata** | Read-only |
| **Contents** | Read-only |
| **Issues** | Read-only |
| **Pull requests** | Read-only |
| **Packages** | Read-only |
| **Repository security advisories** | Read-only |
| **Projects** | Read-only |

### Code scanning permissions

To use the code scanning query type, the following additional permissions are required for both personal access tokens and GitHub Apps:

| Permission | Access level |
|------------|-------------|
| **Code scanning alerts** | Read-only |
| **Security events** | Read-only |

For classic personal access tokens, add the `security_events` scope.

## Verify the connection

After you have added your GitHub connection settings, click **Save & test** to test and save the data source connection. When the connection is successful, you see the message **Data source is working**.

If the connection fails, check the following error messages:

| Error message | Solution |
|---------------|----------|
| "401 Unauthorized. Check your API key/Access token" | Verify your access token is correct and hasn't expired. Ensure it has the required scopes. |
| "404 Not Found. Check the Github Enterprise Server URL" | Verify the GitHub Enterprise Server URL is correct. |
| "Unable to reach the Github Enterprise Server URL from the Grafana server" | Check network connectivity, firewall rules, and proxy settings. For Grafana Cloud, configure [Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/). |

## Provision the data source

You can define the data source in YAML files as part of Grafana's provisioning system. For more information, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

### Personal access token example

```yaml
apiVersion: 1

datasources:
  - name: GitHub
    type: grafana-github-datasource
    jsonData:
      selectedAuthType: personal-access-token
    secureJsonData:
      accessToken: <ACCESS_TOKEN>
```

### GitHub App example

```yaml
apiVersion: 1

datasources:
  - name: GitHub
    type: grafana-github-datasource
    jsonData:
      selectedAuthType: github-app
      appId: <APP_ID>
      installationId: <INSTALLATION_ID>
    secureJsonData:
      privateKey: <PRIVATE_KEY>
```

### GitHub Enterprise Server example

```yaml
apiVersion: 1

datasources:
  - name: GitHub Enterprise
    type: grafana-github-datasource
    jsonData:
      selectedAuthType: personal-access-token
      githubPlan: github-enterprise-server
      githubUrl: https://github.example.com
    secureJsonData:
      accessToken: <ACCESS_TOKEN>
```

## Provision with Terraform

You can provision the GitHub data source using the [Grafana Terraform provider](https://registry.terraform.io/providers/grafana/grafana/latest/docs). For more information, refer to [Provision Grafana with Terraform](https://grafana.com/docs/grafana/latest/administration/infrastructure-as-code/terraform/).

### Personal access token example

```hcl
resource "grafana_data_source" "github" {
  type = "grafana-github-datasource"
  name = "GitHub"

  json_data_encoded = jsonencode({
    selectedAuthType = "personal-access-token"
  })

  secure_json_data_encoded = jsonencode({
    accessToken = var.github_access_token
  })
}
```

### GitHub App example

```hcl
resource "grafana_data_source" "github" {
  type = "grafana-github-datasource"
  name = "GitHub"

  json_data_encoded = jsonencode({
    selectedAuthType = "github-app"
    appId            = var.github_app_id
    installationId   = var.github_installation_id
  })

  secure_json_data_encoded = jsonencode({
    privateKey = var.github_private_key
  })
}
```
