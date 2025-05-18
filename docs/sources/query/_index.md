---
title: Query the GitHub data source
menuTitle: Query
description: Learn how to query GitHub using the GitHub data source plugin
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
weight: 200
---
# Query the GitHub data source

The GitHub data source plugin for Grafana enables you to query and visualize data directly from your GitHub repositories and organizations. With this plugin, you can monitor repository activity, track issues and pull requests, analyze workflow runs, and more from within Grafana.

## Query types

The data source supports the following queries, which you can select from in the query editor under the `Query Type` dropdown:

- [**Commits**](#commits): Retrieve a list of commits for a repository or branch, including commit message, author, and timestamp.
- [**Issues**](#issues): List issues in a repository, using the GitHub query syntax to filter the response.
- [**Contributors**](#contributors): Get a list of contributors to a repository.
- [**Tags**](#tags): List created tags for a repository.
- [**Releases**](#releases): List created releases for a repository.
- [**Pull requests**](#pull-requests): List pull requests for a repository, using the GitHub query syntax to filter the response.
- [**Labels**](#labels): List labels defined in a repository.
- [**Repositories**](#repositories): List repositories for a user or organization.
- [**Milestones**](#milestones): Retrieve milestones for a repository, which can be used to group issues and pull requests.
- [**Packages**](#packages): List packages published in a repository or organization.
- [**Vulnerabilities**](#vulnerabilities): Query security vulnerabilities detected in a repository.
- [**Projects**](#projects): List classic projects associated with a repository or organization.
- [**Stargazers**](#stargazers): Get a list of users who have starred a repository.
- [**Workflows**](#workflows): List GitHub Actions workflows defined in a repository.
- [**Workflow usage**](#workflow-usage): Retrieve usage statistics for a workflow, such as run counts and durations.
- [**Workflow runs**](#workflow-runs): List runs for a specific workflow, including status, conclusion, and timing information.

### Commits

Retrieve a list of commits for a repository or branch, including commit messages, authors, and timestamps. Useful for tracking code changes, deployment activity, or contributor history.

#### Query options

| Name | Description | Required (yes/no) |
| ---- | ----------- | ----------------- |
| Owner | The GitHub user or organization that owns the repository. | Yes |
| Repository | The name of the repository | Yes |
| Ref (Branch/Tag) | The branch or tag to list commits against | Yes |

##### Sample queries
Show all commits to the `main` branch of the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- Ref: `main`

Show all commits against a tag

- Owner: `grafana`
- Repository: `grafana`
- Ref: `v12.0.0`

#### Response

| Name | Description |
| ---- | ----------- |
| id | commit ID |
| author | Name of the commit author |
| author_login | GitHub handle of the commit author |
| author_company | Company name of the commit author |
| committed_at | YYYY-MM-DD HH:MM:SS |
| pushed_at | YYYY-MM-DD HH:MM:SS |
| message | commit message |

### Issues

List issues in a repository using the GitHub query syntax to filter the response. Useful for tracking open bugs, feature requests, or project tasks. 

#### Query options

| Name | Description | Required (yes/no) |
| ---- | ----------- | ----------------- |
| Owner | The GitHub user or organization that owns the repository. | Yes |
| Repository | The name of the repository | Yes |
| Ref (Branch/Tag) | The branch or tag to list commits against | Yes |

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: (Optional) The name of the repository.
- **Query**: (Optional)  A GitHub search query string to filter issues using GitHub's advanced search syntax. This allows you to search by keywords, labels, assignee, author, milestone, state, and more. For details on supported syntax, see [Searching issues and pull requests](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests).
- **Time Field**: The time field to filter the responses on - can be: `CreatedAt`, `ClosedAt` or `UpdatedAt`  

##### Sample queries 
Show all closed issues labeled `type/bug` in the grafana repository.

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:closed label:type/bug`

Show all issues with 'sql expressions' in the title
- Owner: `grafana`
- Repository: `grafana`
- Query: `sql expressions in:title`

#### Response

| Name | Description |
| ---- | ----------- |
| title | Issue title |
| author | GitHub handle of the author |
| author_company | Company name of the commit author |
| repo | Issue repository |
| number | Issue number |
| closed | true / false |
| created_at | YYYY-MM-DD HH:MM:SS |
| closed_at | YYYY-MM-DD HH:MM:SS |
| updated_at | YYYY-MM-DD HH:MM:SS |
| labels | Array of labels i.e. [ "type/bug", "needs more info"] |

{{< admonition type="note" >}}
This query returns a maximum of 1000 results.
{{< /admonition >}}

### Contributors

Get a list of contributors to a repository. 

#### Query options

| Name | Description | Required (yes/no) |
| ---- | ----------- | ----------------- |
| Owner | The GitHub user or organization that owns the repository. | Yes |
| Repository | The name of the repository. | |
| Query | Filter for contributors by name or GitHub handle | No |

##### Sample queries
Show all contributors to the `grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

Search for contributors with `bob` in their name or handle

- Owner: `grafana`
- Repository: `grafana`
- Query: `bob`

#### Response

| Name | Description |
| ---- | ----------- |
| Name | Name of the contributor |
| author | Name of the commit author |
| author_login | GitHub handle of the commit author |
| author_company | Company name of the commit author |
| committed_at | YYYY-MM-DD HH:MM:SS |
| pushed_at | YYYY-MM-DD HH:MM:SS |
| message | commit message |

{{< admonition type="note" >}}
This query returns a maximum of 200 results.
{{< /admonition >}}

### Tags

List created tags for a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

##### Sample queries
Show all tags created for the `grafana` repository within the current selected time range.

- Owner: `grafana`
- Repository: `grafana`

#### Response

| Name | Description |
| ---- | ----------- |
| Name | Name of tag |
| id | Sha for the tag|
| author | Name of the GitHub user who created the tag |
| author_login | GitHub handle of the GitHub user who created the tag |
| author_company | Company name of the GitHub user who created the tag |
| date | YYYY-MM-DD HH:MM:SS |

### Releases

List created releases for a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

##### Sample queries
Show all releases for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository : `grafana`

#### Response

| Name | Description |
| ---- | ----------- |
| Name | Name of release |
| created_by | Name of the GitHub user who created the release |
| is_draft | true / false|
| is_prerelease | true / false|
| tag | Tag name associated with the release |
| url | URL for the tag associated with the release |
| created_at | YYYY-MM-DD HH:MM:SS |
| published_at | YYYY-MM-DD HH:MM:SS |

### Pull requests

List pull requests for a repository, using the GitHub query syntax to filter the response.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: (Optional) The name of the repository.


##### Sample queries
Show all open pull requests authored by `octocat` in the `grafana/grafana` repository.

- Owner : `grafana`
- Repository : `grafana`
- State: `open`
- Author: `octocat`

### Labels

Get all labels defined in a repository, useful for categorizing issues and pull requests.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

##### Sample queries
Show all labels for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

#### Response

| Name | Description |
| color | Hexadecimal number |
| name | Label name | 
| description | Label description |

### Repositories

List repositories for a user or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.

##### Sample queries
Show all repositories for the `grafana` organization.

- Organization: `grafana`

### Milestones

Retrieve milestones for a repository, which can be used to group issues and pull requests.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **State**: (Optional) Filter by milestone state (`open`, `closed`, or `all`).

##### Sample queries
Show all open milestones for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- State: `open`

### Packages

List packages published in a repository or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.
- **Repository**: (Optional) The name of the repository.

##### Sample queries
Show all packages for the `grafana` organization.

- Organization: `grafana`

### Vulnerabilities

Query security vulnerabilities detected in a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

##### Sample queries
Show all security advisories for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Projects

List classic projects associated with a repository or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.
- **Repository**: (Optional) The name of the repository.

##### Sample queries
Show all projects for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Stargazers

Get a list of users who have starred a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

##### Sample queries
Show all stargazers for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Workflows

List GitHub Actions workflows defined in a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

##### Sample queries
Show all workflows for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Workflow usage

Retrieve usage statistics for a workflow, such as run counts and durations.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **Workflow**: The workflow to get usage for.

##### Sample queries
Show usage statistics for the `CI` workflow in the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- Workflow: `CI`


### Workflow runs

List runs for a specific workflow, including status, conclusion, and timing information.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **Workflow**: The workflow to list runs for.
- **Branch**: (Optional) Filter by branch.
- **Status**: (Optional) Filter by status (`queued`, `in_progress`, `completed`).
- **Conclusion**: (Optional) Filter by conclusion (`success`, `failure`, etc.).
- **Created**: (Optional) Filter by creation date.

##### Sample queries
Show all completed runs for the `CI` workflow on the `main` branch in the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- Workflow: `CI`
- Branch: `main`
- Status: `completed`
