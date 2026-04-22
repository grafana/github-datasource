---
title: Create a GitHub access token
menuTitle: Create an access token
description: Create a GitHub access token
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
weight: 102
---

# Create a GitHub access token

You will need either a `GitHub App` or a `Personal Access Token` to use this plugin.

## Creating a personal access token (classic)

This is an example when you want to use the personal access token (classic). \
Read more about [personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).

1. Login to your GitHub account.
1. Navigate to [Personal access tokens](https://github.com/settings/tokens) and click **Generate new token**.
1. Select the **personal access token (classic)**.
1. Define the permissions which you want to allow.
1. Click **Generate Token**.

## Creating a fine-grained personal access token

This is an example when you want to use the fine-grained personal access token. \
Read more about [fine-grained personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token).

1. Login to your GitHub account.
1. Navigate to [Personal access tokens](https://github.com/settings/tokens?type=beta) and click **Generate new token**.
1. Provide a name for the token.
1. Provide necessary permissions which you want to allow. Ensure you are providing `read-only` permissions.
1. Click **Generate token**.

## Using GitHub App Authentication

You can also authenticate using a GitHub App instead of a personal access token. This method allows for better security and fine-grained access to resources.

1. Register a new GitHub App by following the instructions in the [GitHub App documentation](https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app).
1. After registering the App, generate a private key for authentication.
1. Note down the App ID assigned to your GitHub App.
1. [Install the GitHub App](https://docs.github.com/en/apps/using-github-apps/installing-your-own-github-app) on your GitHub account or organization.
1. Note the installation ID after completing the installation.
1. In Grafana's data source settings, provide the **app id**, **installation id**, and **private key** in the appropriate fields.

> **Where to find your installation id?** \
> Head over to Settings > Installed GitHub Apps > Configure. The installation ID can be found at the end of the URL `https://github.com/settings/installations/<installation_id>`.

## Permissions

You will need to define the access permissions for your **personal access token** in order to allow it to access the data.

The following lists include the required permissions for the access token:

`public_repo`, `repo:status`, `repo_deployment`, `read:packages`, `read:user`, `user:email`, `read:org`, `read:project`, `repo (For full control of private repositories)`

You will need to define the access permissions for your **GitHub App** in order to allow it to access the data.

**Repositories:**

`metadata: read-only`, `contents: read-only`, `issues: read-only`, `pull requests: read-only`, `packages: read-only`, `repository security advisories: read-only`, `projects: read-only`

**Code scan**

`code scanning alerts: read-only, security_events: read-only`
