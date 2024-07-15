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

You will need a **personal access token** to use the plugin. Below is a table that indicates what minimum requirements must be matched before the plugin can be used.

Options:

| Setting               | Required |  Details                                              |
| --------------------- | -------- |-------------------------------------------------------|
| Access token          | true     | This is required to allow plugin to connect to GitHub |
| Default Organization  | false    | Only if you want to set a specfic organization as a default selection while creating Dashboards    |
| Default Repository    | true     | A repository is required to access the data           |
| GitHub Enterprise URL | false    | Only if you are using GitHub Enterprise account            |

Read more about [creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).

To create a new personal access token, navigate to [Personal access tokens](https://github.com/settings/tokens) and press Generate new token.

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
