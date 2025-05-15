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

- [**Commits**](#commits): Retrieve a list of commits for a repository or branch, including commit messages, authors, and timestamps.
- [**Issues**](#issues): List issues in a repository, with options to filter by state, assignee, and labels.
- [**Contributors**](#contributors): Get a list of contributors to a repository, including their contribution counts.
- [**Tags**](#tags): List tags for a repository, often used to identify release points.
- [**Releases**](#releases): Retrieve information about releases published in a repository.
- [**Pull requests**](#pull-requests): List pull requests for a repository, with filtering by state, author, and labels.
- [**Labels**](#labels): Get all labels defined in a repository, useful for categorizing issues and pull requests.
- [**Repositories**](#repositories): List repositories for a user or organization.
- [**Milestones**](#milestones): Retrieve milestones for a repository, which can be used to group issues and pull requests.
- [**Packages**](#packages): List packages published in a repository or organization.
- [**Vulnerabilities**](#vulnerabilities): Query security vulnerabilities detected in a repository.
- [**Projects**](#projects): List classic projects associated with a repository or organization.
- [**Stargazers**](#stargazers): Get a list of users who have starred a repository.
- [**Workflows**](#workflows): List GitHub Actions workflows defined in a repository.
- [**Workflow usage**](#workflow-usage): Retrieve usage statistics for a workflow, such as run counts and durations.
- [**Workflow runs**](#workflow-runs): List runs for a specific workflow, including status, conclusion, and timing information.

---

## Examples

### Commits

Retrieve a list of commits for a repository or branch, including commit messages, authors, and timestamps. Useful for tracking code changes, deployment activity, or contributor history.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **Ref (Branch/Tag)**: The branch to filter commits by.

**Sample queries:**  
Show all commits to the `main` branch of the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- Ref: `main`

Show all commits against a tag

- Owner: `grafana`
- Repository: `grafana`
- Ref: `v12.0.0`


---

### Issues

List issues in a repository, with options to filter by state, assignee, and labels. Useful for tracking open bugs, feature requests, or project tasks.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **Query**: A GitHub search query string to filter issues using GitHub's advanced search syntax. This allows you to search by keywords, labels, assignee, author, milestone, state, and more. For details on supported syntax, see [Searching issues and pull requests](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests).
- **Time Field**: The time field to filter the responses on - can be: `CreatedAt`, `ClosedAt` or `UpdatedAt`  

**Sample queries:**  
Show all closed issues labeled `type/bug`

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:closed label:type/bug`

Show all issues with 'sql expressions' in the title
- Owner: `grafana`
- Repository: `grafana`
- Query: `sql expressions in:title`
---

### Contributors

Get a list of contributors to a repository, including their contribution counts. Useful for understanding project activity and top contributors.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **Query (optional)**: Filter for contributors by name or GitHub handle

**Sample queries:**  
Show all contributors to the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

Search for contributors with `bob` in their name or handle

- Owner: `grafana`
- Repository: `grafana`
- Query: `bob`


---

### Tags

List tags for a repository, often used to identify release points.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all tags for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

---

### Releases

Retrieve information about releases published in a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all releases for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

---

### Pull requests

List pull requests for a repository, with filtering by state, author, and labels.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **State**: Filter by PR state (`open`, `closed`, or `all`).
- **Author**: (Optional) Filter by author.
- **Labels**: (Optional) Filter by labels.
- **Base**: (Optional) Filter by base branch.

**Sample query:**  
Show all open pull requests authored by `octocat` in the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- State: `open`
- Author: `octocat`

### Labels

Get all labels defined in a repository, useful for categorizing issues and pull requests.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all labels for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Repositories

List repositories for a user or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.

**Sample query:**  
Show all repositories for the `grafana` organization.

- Organization: `grafana`

### Milestones

Retrieve milestones for a repository, which can be used to group issues and pull requests.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **State**: (Optional) Filter by milestone state (`open`, `closed`, or `all`).

**Sample query:**  
Show all open milestones for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- State: `open`

### Packages

List packages published in a repository or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.
- **Repository**: (Optional) The name of the repository.

**Sample query:**  
Show all packages for the `grafana` organization.

- Organization: `grafana`

### Vulnerabilities

Query security vulnerabilities detected in a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all security advisories for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Projects

List classic projects associated with a repository or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.
- **Repository**: (Optional) The name of the repository.

**Sample query:**  
Show all projects for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Stargazers

Get a list of users who have starred a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all stargazers for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Workflows

List GitHub Actions workflows defined in a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all workflows for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

### Workflow usage

Retrieve usage statistics for a workflow, such as run counts and durations.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **Workflow**: The workflow to get usage for.

**Sample query:**  
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

**Sample query:**  
Show all completed runs for the `CI` workflow on the `main` branch in the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- Workflow: `CI`
- Branch: `main`
- Status: `completed`
