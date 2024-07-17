---
title: Create a GitHub personal access token
menuTitle: Create a personal access token
description: Create a GitHub personal access token
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

# Create a GitHub personal access token

You will need a _personal access token_ to use the plugin. GitHub currently supports two types of personal access tokens:

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

### Permissions

You will need to define the access permissions for your token in order to allow it to access the data.

The following lists include the required permissions for the access token:

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
