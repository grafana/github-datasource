---
slug: '/configuration'
title: 'Configuration'
menuTitle: Configuration
description: Configuration of the GitHub data source plugin
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

# Configuration

You will need a **GitHub Access Token** to use the plugin. Below is a table that indicates what minimum requirements must be matched before the plugin can be used.

Options:

| Setting               | Required |
| --------------------- | -------- |
| Access token          | true     |
| Default Organization  | false    |
| Default Repository    | true     |
| GitHub Enterprise URL | false    |

To create a new Access Token, navigate to [Personal Access Tokens](https://github.com/settings/tokens) and press Generate new token.

## Access token permissions

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

An extra setting is required for private repositories

- `repo (Full control of private repositories)`
