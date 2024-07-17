---
title: 'Configure Token'
menuTitle: Configure Token
description: Configuring the GitHub personal access token
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

# Configure Token

You will need a **personal access token** to use the plugin. GitHub currently supports two types of personal access tokens:

1. fine-grained personal access tokens
1. personal access tokens (classic)

Read more about [personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).

The Grafana GitHub data source plugin works with both. Below is a table that indicates what minimum requirements must be matched before the plugin can be used.

Options:

| Setting               | Required | Description                                                                                     |
| --------------------- | -------- | ----------------------------------------------------------------------------------------------- |
| Access token          | true     | This is required to allow plugin to connect to GitHub                                           |
| GitHub Enterprise URL | false    | Only if you are using GitHub Enterprise account                                                 |

## Creating a personal access token (classic)

This is an example when you want to use the personal access token (classic).

1. Login to your GitHub account.
1. Navigate to [Personal access tokens](https://github.com/settings/tokens) and click **Generate new token**.
1. Select the **personal access token (classic)**.
1. Define the permissions which you want to allow.
1. Click **Generate Token**.

### personal access token permissions

You will need to define the access permissions for your token in order to allow it to access the data.

Here is a list of the required minimum permissions defined below:

For all repositories:

- `public_repo`
- `repo:status`
- `repo_deployment`
- `read:packages`
- `read:user`
- `user:email`

For GitHub projects:

- `read:org`
- `read:project`

An extra setting is required for private repositories:

- `repo (Full control of private repositories)`
