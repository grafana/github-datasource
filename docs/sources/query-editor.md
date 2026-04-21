---
aliases:
  - ./query/
title: GitHub query editor
menuTitle: Query editor
description: Learn how to query GitHub data using the GitHub data source plugin query editor
keywords:
  - data source
  - github
  - query
  - commits
  - issues
  - pull requests
  - workflows
  - deployments
  - code scanning
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 200
review_date: "2026-03-11"
---

# GitHub query editor

This document describes the query types and response fields available in the GitHub data source query editor. Most query types use common **Owner** and **Repository** fields, and some support additional filtering through the **Query** field using [GitHub search syntax](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests).

## Query types

Select a query type from the **Query Type** drop-down in the query editor:

- [**Code scanning**](#code-scanning): Query code scanning alerts for a repository or organization.
- [**Commit files**](#commit-files): List files changed in a specific commit.
- [**Commits**](#commits): Retrieve a list of commits for a branch or ref within a repository, including commit message, author, and timestamp.
- [**Contributors**](#contributors): Get a list of contributors to a repository.
- [**Deployments**](#deployments): List deployments for a repository, including environment, ref, and task information.
- [**Issues**](#issues): List issues in a repository, using the GitHub query syntax to filter the response.
- [**Labels**](#labels): List labels defined in a repository.
- [**Milestones**](#milestones): Retrieve milestones for a repository, which can be used to group issues and pull requests.
- [**Packages**](#packages): List packages published from a repository in an organization.
- [**Projects**](#projects): List projects associated with a user or organization.
- [**Pull request files**](#pull-request-files): List files changed in a specific pull request.
- [**Pull requests**](#pull-requests): List pull requests for a repository, using the GitHub query syntax to filter the response.
- [**Pull request reviews**](#pull-request-reviews): List reviews for pull requests in a repository.
- [**Releases**](#releases): List created releases for a repository.
- [**Repositories**](#repositories): List repositories for a user or organization.
- [**Stargazers**](#stargazers): Get a list of users who have starred a repository, including the ability to plot a total count over time.
- [**Tags**](#tags): List created tags for a repository.
- [**Vulnerabilities**](#vulnerabilities): Query security vulnerabilities detected in a repository.
- [**Workflows**](#workflows): List GitHub Actions workflows defined in a repository.
- [**Workflow runs**](#workflow-runs): List runs for a specific workflow, including status, conclusion, and timing information.
- [**Workflow usage**](#workflow-usage): Retrieve usage statistics for a workflow, such as run counts and durations.

### Code scanning

Query code scanning alerts for a repository or organization. Useful for tracking security issues detected by GitHub's code scanning tools.

{{< admonition type="note" >}}
Code scanning requires additional permissions. Refer to [Code scanning permissions](https://grafana.com/docs/plugins/grafana-github-datasource/latest/configure/#code-scanning-permissions).
{{< /admonition >}}

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository. Leave empty to query alerts at the organization level. | No |
| State | Filter by alert state: `open`, `closed`, `dismissed`, or `fixed` | No |
| Ref | The Git reference to filter results. Format as `refs/heads/<BRANCH_NAME>` or `<BRANCH_NAME>`. For pull requests, use `refs/pull/<NUMBER>/merge`. | No |

##### Sample queries

Show all open code scanning alerts for the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`
- State: `open`

Show all code scanning alerts across the `grafana` organization:

- Owner: `grafana`
- Repository: _(empty)_
- State: _(empty)_

#### Response

| Name | Description |
|------|-------------|
| number | Alert number |
| created_at | When the alert was created: YYYY-MM-DD HH:MM:SS |
| updated_at | When the alert was last updated: YYYY-MM-DD HH:MM:SS |
| dismissed_at | When the alert was dismissed, if applicable: YYYY-MM-DD HH:MM:SS |
| url | URL to the alert in the GitHub web UI |
| state | Alert state: `open`, `closed`, `dismissed`, or `fixed` |
| dismissed_by | GitHub handle of the user who dismissed the alert |
| dismissed_reason | Reason the alert was dismissed |
| dismissed_comment | Comment provided when the alert was dismissed |
| rule_id | Identifier for the rule that triggered the alert |
| rule_severity | Severity level of the rule |
| rule_security_severity_level | Security severity level of the rule |
| rule_description | Short description of the rule |
| rule_full_description | Full description of the rule |
| rule_tags | Comma-separated list of tags associated with the rule |
| rule_help | Help text for the rule |
| tool_name | Name of the code scanning tool |
| tool_version | Version of the code scanning tool |
| tool_guid | GUID of the code scanning tool |

### Commits

Retrieve a list of commits for a branch or ref within a repository, including commit message, author, and timestamp. Useful for tracking code changes, deployment activity, or contributor history.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |
| Ref (Branch/Tag) | The branch or tag to list commits against | Yes |
| Include files | Include file-level change details for each commit. For details, refer to [Response with Include files enabled](#response-with-include-files-enabled). | No |

##### Sample queries

Show all commits to the `main` branch of the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`
- Ref: `main`

Show all commits against a tag:

- Owner: `grafana`
- Repository: `grafana`
- Ref: `v12.0.0`

#### Response

| Name | Description |
|------|-------------|
| id | Commit ID |
| author | Name of the commit author |
| author_login | GitHub handle of the commit author |
| author_email | Email address of the commit author |
| author_company | Company name of the commit author |
| committed_at | When the change was committed: YYYY-MM-DD HH:MM:SS |
| pushed_at | When the commit was pushed: YYYY-MM-DD HH:MM:SS |
| message | The commit message |

#### Response with Include files enabled

{{< admonition type="warning" >}}
This option makes one additional API call per commit. Avoid enabling it over large time ranges to prevent rate limiting.
{{< /admonition >}}

When **Include files** is enabled, the response returns one row for each changed file in each commit:

| Name | Description |
|------|-------------|
| id | Commit ID |
| author | Name of the commit author |
| author_login | GitHub handle of the commit author |
| author_email | Email address of the commit author |
| author_company | Company name of the commit author |
| committed_at | When the change was committed: YYYY-MM-DD HH:MM:SS |
| pushed_at | When the commit was pushed: YYYY-MM-DD HH:MM:SS |
| message | The commit message |
| file_path | Path of the changed file |
| file_additions | Number of lines added in this file |
| file_deletions | Number of lines deleted in this file |
| file_changes | Total number of line changes in this file |
| file_status | Change type: `added`, `modified`, `renamed`, etc. |
| previous_filename | Original path for renamed files |

### Commit files

List files changed in a specific commit.

{{< admonition type="note" >}}
The GitHub API returns at most 300 files for a single commit.
{{< /admonition >}}

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |
| Commit SHA | The SHA of the target commit | Yes |

#### Response

The response includes one row per changed file:

| Name | Description |
|------|-------------|
| path | Path of the changed file |
| additions | Number of lines added |
| deletions | Number of lines deleted |
| changes | Total number of line changes |
| status | Change type: `added`, `modified`, `renamed`, etc. |
| previous_filename | Original path for renamed files |

### Contributors

Get a list of contributors to an organization or repository.

{{< admonition type="note" >}}
This query returns a maximum of 200 results.
{{< /admonition >}}

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |
| Query | Filter for contributors by name or GitHub handle | No |

##### Sample queries

Show all contributors to the `grafana` repository:

- Owner: `grafana`
- Repository: `grafana`

Search for contributors with `bob` in their name or handle:

- Owner: `grafana`
- Repository: `grafana`
- Query: `bob`

#### Response

| Name | Description |
|------|-------------|
| name | Name of the contributor |
| author | Name of the commit author |
| author_login | GitHub handle of the commit author |
| author_company | Company name of the commit author |
| committed_at | When the commit was made: YYYY-MM-DD HH:MM:SS |
| pushed_at | When the commit was pushed: YYYY-MM-DD HH:MM:SS |
| message | Commit message |

### Deployments

List deployments for a repository, including environment, ref, and task information. Deployments track deployment requests for specific refs (branches, tags, or SHAs) to different environments.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | GitHub user or organization that owns the repository | Yes |
| Repository | Name of the repository | Yes |
| SHA | Filter deployments by the SHA recorded at creation time | No |
| Ref | Filter by ref name (branch, tag, or SHA) | No |
| Task | Filter by task name (for example, `deploy` or `deploy:migrations`) | No |
| Environment | Filter by environment name (for example, `production`, `staging`, or `qa`) | No |

##### Sample queries

Show all deployments for the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`

Show all deployments to the production environment:

- Owner: `grafana`
- Repository: `grafana`
- Environment: `production`

Show all deployments for a specific branch:

- Owner: `grafana`
- Repository: `grafana`
- Ref: `main`

#### Response

| Name | Description |
|------|-------------|
| id | Unique identifier for the deployment |
| sha | SHA of the commit that was deployed |
| ref | Ref (branch, tag, or SHA) that was deployed |
| task | Task name (for example, `deploy` or `deploy:migrations`) |
| environment | Environment name (for example, `production` or `staging`) |
| description | Description of the deployment |
| creator | GitHub handle of the user who created the deployment |
| created_at | When the deployment was created: YYYY-MM-DD HH:MM:SS |
| updated_at | When the deployment was last updated: YYYY-MM-DD HH:MM:SS |
| url | API URL for the deployment |
| statuses_url | API URL for the deployment statuses |

### Issues

List issues in a repository using the GitHub query syntax to filter the response. Useful for tracking open bugs, feature requests, or project tasks.

{{< admonition type="note" >}}
This query returns a maximum of 1000 results.
{{< /admonition >}}

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | A GitHub user or organization | Yes |
| Repository | The name of a repository | No |
| Query | Use GitHub's [query syntax](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests) to filter results | No |
| Time field | The time field to filter the responses on: `CreatedAt`, `ClosedAt`, or `UpdatedAt` | Yes |

##### Sample queries

Show all closed issues labeled `type/bug` in the Grafana repository:

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:closed label:type/bug`

Show all issues with "sql expressions" in the title:

- Owner: `grafana`
- Repository: `grafana`
- Query: `sql expressions in:title`

Show open issues with no assignees (useful for identifying unowned work):

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:open no:assignee`

Show issues created in the last 7 days using the `$__toDay` macro:

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:open created:>=$__toDay(-7)`
- Time field: `CreatedAt`

#### Response

| Name | Description |
|------|-------------|
| title | Issue title |
| author | GitHub handle of the author |
| author_company | Company name of the commit author |
| repo | Issue repository |
| number | Issue number |
| state | Issue state: `OPEN` or `CLOSED` |
| closed | `true` or `false` |
| created_at | When the issue was created: YYYY-MM-DD HH:MM:SS |
| closed_at | When the issue was closed: YYYY-MM-DD HH:MM:SS |
| updated_at | When the issue was last updated: YYYY-MM-DD HH:MM:SS |
| labels | Array of labels, for example: `["type/bug", "needs more info"]` |
| assignees | Array of assignees, for example: `["user1", "user2"]` |
| milestone | Milestone title associated with the issue, if applicable |

### Labels

Get all labels defined in a repository, useful for categorizing issues and pull requests.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |
| Query | Filter on text in name and description for labels | No |

##### Sample queries

Show all labels for the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`

#### Response

| Name | Description |
|------|-------------|
| name | Label name |
| color | Hexadecimal number |
| description | Label description |

### Milestones

Retrieve milestones for a repository, which can be used to group issues and pull requests.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |
| Query | Filter on text in the milestone title | No |

##### Sample queries

Show all milestones for the `grafana/grafana` repository for v11 of Grafana:

- Owner: `grafana`
- Repository: `grafana`
- Query: `11.`

#### Response

| Name | Description |
|------|-------------|
| title | Milestone title |
| author | GitHub handle of the user who created the milestone |
| closed | Whether the milestone is closed: `true` or `false` |
| state | One of `OPEN` or `CLOSED` |
| created_at | When the milestone was created: YYYY-MM-DD HH:MM:SS |
| closed_at | When the milestone was closed: YYYY-MM-DD HH:MM:SS |
| due_at | When the milestone is due by: YYYY-MM-DD HH:MM:SS |

{{< admonition type="note" >}}
Milestone titles can be anything and are therefore parsed as a string. This means sorting by title may appear incorrect if you have numeric milestones, for example: `12.0.0`. [Transformations](https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/transform-data/) can be used to change the data type in this scenario.
{{< /admonition >}}

### Packages

List packages published from a repository in an organization.

{{< admonition type="note" >}}
This query doesn't support querying npm, RubyGems, or NuGet packages.
{{< /admonition >}}

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |
| Package type | One of: `MAVEN`, `DOCKER`, `DEBIAN`, or `PYPI` | Yes |
| Names | Filter for packages using a comma-separated list of names | No |

##### Sample queries

Show all packages uploaded to the `grafana` organization:

- Organization: `grafana`

#### Response

| Name | Description |
|------|-------------|
| name | Package name |
| platform | Platform or registry where the package is published |
| version | Package version |
| type | Package type (for example, `MAVEN`, `DOCKER`, `DEBIAN`, `PYPI`) |
| prerelease | Whether the package version is a prerelease: `true` or `false` |
| downloads | Number of downloads for the package version |

### Projects

List projects associated with a user or organization.

{{< admonition type="note" >}}
This query returns a maximum of 200 results.
{{< /admonition >}}

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Project Owner | One of `Organization` or `User` | Yes |
| Organization | Organization for the project (shown when Organization is selected) | Yes |
| User | User for the project (shown when User is selected) | Yes |
| Project Number | Enter a specific project number to query for associated items | No |
| Filter | Add key-value filters based on the fields for project items (shown if Project Number is specified) | No |

##### Sample queries

Show all projects for the `grafana` organization:

- Project Owner: `Organization`
- Organization: `grafana`

Show all pull requests for the "Dashboards" project in the Grafana organization:

- Project Owner: `Organization`
- Organization: `grafana`
- Project Number: `202`
- Filter: `type equal PULL_REQUEST`

#### Response

The response format depends on whether a project number is specified.

##### Without a project number

| Name | Description |
|------|-------------|
| number | The project number |
| title | Title of the project |
| url | URL for the project |
| closed | Whether the project has been closed: `true` or `false` |
| public | Whether the project is public: `true` or `false` |
| closed_at | When the project was closed: YYYY-MM-DD HH:MM:SS |
| updated_at | When the project was last updated |
| created_at | When the project was created |
| short_description | The description of the project |

##### With a project number

{{< admonition type="note" >}}
GitHub Projects allow for customization of default fields and custom fields to be added. Therefore, the response can vary significantly between projects.
{{< /admonition >}}

| Name | Description |
|------|-------------|
| name | Name of the project item (issue or pull request) |
| id | Unique identifier for the project item |
| type | Type of the item (such as `ISSUE`, `PULL_REQUEST`) |
| status | Status of the item (such as "In development", "Shipped") -- this can be configured on the project |
| labels | Comma-separated list of labels assigned to the item |
| assignees | Comma-separated list of users assigned to the item |
| reviewers | Comma-separated list of reviewers (for pull requests) |
| repository | Name of the repository the item belongs to |
| milestone | Milestone associated with the item |
| priority | Priority value or label |
| archived | Whether the item is archived: `true` or `false` |
| created_at | When the item was created: YYYY-MM-DD HH:MM:SS |
| updated_at | When the item was last updated: YYYY-MM-DD HH:MM:SS |
| closed_at | When the item was closed, if applicable: YYYY-MM-DD HH:MM:SS |
| (custom fields) | Any custom defined fields are also returned alongside their values |

### Pull requests

List pull requests for a repository, using the GitHub query syntax to filter the response.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | A GitHub user or organization | Yes |
| Repository | The name of a repository | No |
| Query | Use GitHub's [query syntax](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests) to filter results | No |
| Time field | The time field to filter the responses on: `CreatedAt`, `ClosedAt`, `MergedAt`, `UpdatedAt`, or `None` | Yes |

##### Sample queries

Show all open pull requests authored by Renovate in the `grafana/plugin-tools` repository:

- Owner: `grafana`
- Repository: `plugin-tools`
- Query: `is:open author:app/renovate`
- Time field: `None`

Show pull requests merged in the selected time range:

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:merged`
- Time field: `MergedAt`

Show open draft pull requests (useful for tracking work in progress):

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:open draft:true`
- Time field: `None`

Show open pull requests with a specific label:

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:open label:priority/high`
- Time field: `None`

#### Response

| Name | Description |
|------|-------------|
| number | Pull request number |
| title | Pull request title |
| url | URL to the pull request |
| additions | Total number of lines of code added or altered in the pull request |
| deletions | Total number of lines of code removed or altered in the pull request |
| repository | Repository for the pull request |
| state | Can be `open`, `closed`, or `merged` |
| author_name | Name of the GitHub user who created the pull request |
| author_login | GitHub handle of the GitHub user who created the pull request |
| author_email | Email address of the GitHub user who created the pull request |
| author_company | Company name of the GitHub user who created the pull request |
| closed | Whether the pull request is closed: `true` or `false` |
| is_draft | Whether the pull request is in draft: `true` or `false` |
| locked | Whether the pull request has been locked: `true` or `false` |
| merged | Whether the pull request has been merged |
| mergeable | Whether the pull request can be automatically merged: `MERGEABLE`, `CONFLICTING`, or `UNKNOWN` |
| closed_at | When the pull request was closed: YYYY-MM-DD HH:MM:SS |
| merged_at | When the pull request was merged: YYYY-MM-DD HH:MM:SS |
| merged_by_name | Name of the GitHub user who merged the pull request |
| merged_by_login | GitHub handle of the GitHub user who merged the pull request |
| merged_by_email | Email address of the GitHub user who merged the pull request |
| merged_by_company | Company name of the GitHub user who merged the pull request |
| updated_at | When the pull request was last updated: YYYY-MM-DD HH:MM:SS |
| created_at | When the pull request was created: YYYY-MM-DD HH:MM:SS |
| open_time | Duration in seconds the pull request has been open |
| labels | Array of labels assigned to the pull request, for example: `["bug", "priority/high"]` |

### Pull request files

List all files changed in a specific pull request.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |
| Pull Request Number | The number of the pull request | Yes |

#### Response

The response includes one row per changed file:

| Name | Description |
|------|-------------|
| path | Path of the changed file |
| additions | Number of lines added |
| deletions | Number of lines deleted |
| changes | Total number of line changes |
| status | Change type: `added`, `modified`, `renamed`, etc. |
| previous_filename | Original path for renamed files |

### Pull request reviews

List reviews for pull requests in a repository. This query returns one row per review, with the associated pull request details included. Useful for tracking review activity, reviewer workload, and review turnaround time.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | A GitHub user or organization | Yes |
| Repository | The name of a repository | No |
| Query | Use GitHub's [query syntax](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests) to filter results | No |
| Time field | The time field to filter the responses on: `CreatedAt`, `ClosedAt`, `MergedAt`, `UpdatedAt`, or `None` | Yes |

##### Sample queries

Show all reviews for merged pull requests in the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:merged`
- Time field: `MergedAt`

Show reviews for open pull requests (useful for tracking pending review workload):

- Owner: `grafana`
- Repository: `grafana`
- Query: `is:open`
- Time field: `None`

#### Response

| Name | Description |
|------|-------------|
| pull_request_number | Pull request number |
| pull_request_title | Pull request title |
| pull_request_state | Pull request state: `OPEN`, `CLOSED`, or `MERGED` |
| pull_request_url | URL to the pull request |
| pull_request_author_name | Name of the pull request author |
| pull_request_author_login | GitHub handle of the pull request author |
| pull_request_author_email | Email of the pull request author |
| pull_request_author_company | Company of the pull request author |
| repository | Repository name with owner (for example, `grafana/grafana`) |
| review_author_name | Name of the reviewer |
| review_author_login | GitHub handle of the reviewer |
| review_author_email | Email of the reviewer |
| review_author_company | Company of the reviewer |
| review_url | URL to the review |
| review_state | Review state: `APPROVED`, `CHANGES_REQUESTED`, `COMMENTED`, `DISMISSED`, or `PENDING` |
| review_comment_count | Number of comments in the review |
| review_updated_at | When the review was last updated: YYYY-MM-DD HH:MM:SS |
| review_created_at | When the review was created: YYYY-MM-DD HH:MM:SS |

### Releases

List created releases for a repository.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |

##### Sample queries

Show all releases for the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`

#### Response

| Name | Description |
|------|-------------|
| name | Name of release |
| created_by | Name of the GitHub user who created the release |
| is_draft | Whether the release is a draft release: `true` or `false` |
| is_prerelease | Whether the release is a pre-release: `true` or `false` |
| tag | Tag name associated with the release |
| url | URL for the tag associated with the release |
| created_at | When the release was created: YYYY-MM-DD HH:MM:SS |
| published_at | When the release was published: YYYY-MM-DD HH:MM:SS |

### Repositories

List repositories for a user or organization.

{{< admonition type="note" >}}
This query returns a maximum of 1000 results.
{{< /admonition >}}

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | A GitHub user or organization | Yes |
| Repository | Filter on the name of the repository | No |

##### Sample queries

Show all repositories for the `grafana` organization:

- Organization: `grafana`

#### Response

| Name | Description |
|------|-------------|
| name | Name of the repository |
| owner | Organization or user who owns the repository |
| name_with_owner | Returns the owner and repository name in the format `<OWNER>/<REPOSITORY>`, for example: `grafana/loki` |
| url | URL for the repository |
| forks | The number of forks for a repository |
| is_fork | Whether the repository is a fork of another repository: `true` or `false` |
| is_mirror | Whether the repository is a mirror of another repository: `true` or `false` |
| is_private | Whether the repository is private: `true` or `false` |
| created_at | When the repository was created: YYYY-MM-DD HH:MM:SS |

### Stargazers

Get a list of users who have starred a repository, including the ability to plot a total count over time.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |

##### Sample queries

Show all stargazers for the `grafana/grafana` repository within the current time range:

- Owner: `grafana`
- Repository: `grafana`

#### Response

| Name | Description |
|------|-------------|
| starred_at | When the user starred the repository: YYYY-MM-DD HH:MM:SS |
| star_count | Current total of stars for the repository at the time of the event |
| id | `node_id` -- a unique identifier for the GitHub user which can be used in GitHub's GraphQL API |
| login | GitHub handle of the user who starred the repository |
| git_name | Name of the GitHub user who starred the repository |
| company | Company name of the GitHub user who starred the repository |
| email | Email address of the GitHub user who starred the repository |
| url | URL to the GitHub profile for the user who starred the repository |

### Tags

List created tags for a repository.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |

##### Sample queries

Show all tags created for the `grafana/grafana` repository within the currently selected time range:

- Owner: `grafana`
- Repository: `grafana`

#### Response

| Name | Description |
|------|-------------|
| name | Name of the tag |
| id | SHA for the tag |
| author | Name of the user who created the tag |
| author_login | GitHub handle of the user who created the tag |
| author_company | Company name of the user who created the tag |
| date | When the tag was created: YYYY-MM-DD HH:MM:SS |

### Vulnerabilities

Query security vulnerabilities detected in a repository.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | The GitHub user or organization that owns the repository | Yes |
| Repository | The name of the repository | Yes |

##### Sample queries

Show all security advisories for the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`

#### Response

| Name | Description |
|------|-------------|
| value | Custom field which allows for counting or aggregation, always returns `1` |
| created_at | When the vulnerability alert was created: YYYY-MM-DD HH:MM:SS |
| dismissed_at | When the vulnerability alert was dismissed, if applicable: YYYY-MM-DD HH:MM:SS |
| dismissed_reason | Reason the vulnerability alert was dismissed (for example, false positive, won't fix), if applicable |
| withdrawn_at | When the advisory was withdrawn, if applicable: YYYY-MM-DD HH:MM:SS |
| packageName | Name of the affected package |
| advisoryDescription | Description of the vulnerability or advisory |
| firstPatchedVersion | The first version of the package where the vulnerability is fixed |
| vulnerableVersionRange | The range of package versions affected by the vulnerability |
| cvssScore | CVSS (Common Vulnerability Scoring System) score for the vulnerability |
| cvssVector | CVSS vector string describing the scoring metrics |
| permalink | URL to the GitHub Security Advisory or alert |
| severity | Severity level: `LOW`, `MODERATE`, `HIGH`, or `CRITICAL` |
| state | State of the vulnerability alert: `OPEN`, `FIXED`, or `DISMISSED` |

### Workflows

List GitHub Actions workflows defined in a repository.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | GitHub user or organization that owns the repository | Yes |
| Repository | Name of the repository | Yes |
| Time field | The time field to filter the responses on: `None` (returns all workflows), `CreatedAt`, or `UpdatedAt` | Yes |

##### Sample queries

Show all workflows in the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`
- Time field: `None`

Show all workflows created within the `grafana/grafana` repository within the current time range:

- Owner: `grafana`
- Repository: `grafana`
- Time field: `CreatedAt`

#### Response

| Name | Description |
|------|-------------|
| id | Unique identifier for the workflow |
| name | Name of the workflow |
| path | Path to the workflow YAML file in the repository |
| state | State of the workflow: `active`, `deleted`, `disabled_fork`, `disabled_inactivity`, or `disabled_manually` |
| created_at | When the workflow was created: YYYY-MM-DD HH:MM:SS |
| updated_at | When the workflow was last updated: YYYY-MM-DD HH:MM:SS |
| url | API URL for the workflow |
| html_url | URL to the workflow file in the repository |
| badge_url | URL to the workflow status badge |

### Workflow runs

List runs for a specific workflow, including status, conclusion, and timing information.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | GitHub user or organization that owns the repository | Yes |
| Repository | Name of the repository | Yes |
| Workflow | The workflow ID or file name. Use `id` or the filename from `path` from [Workflows](#workflows) queries. | Yes |
| Branch | The head branch to filter on | No |

##### Sample queries

Show all completed runs for the `Levitate` workflow in the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`
- Workflow: `detect-breaking-changes-levitate.yml`

Show workflow runs for a specific branch (useful for monitoring release branch CI):

- Owner: `grafana`
- Repository: `grafana`
- Workflow: `ci.yml`
- Branch: `release/v12.0.x`

#### Response

| Name | Description |
|------|-------------|
| id | Unique identifier for the workflow run |
| name | Name of the workflow or workflow job |
| head_branch | Name of the branch the workflow run was triggered on |
| head_sha | Commit SHA that triggered the workflow run |
| created_at | When the workflow run was created: YYYY-MM-DD HH:MM:SS |
| updated_at | When the workflow run was last updated: YYYY-MM-DD HH:MM:SS |
| html_url | URL to the workflow run in the GitHub web UI |
| url | API URL for the workflow run |
| status | Current status: `queued`, `in_progress`, `completed`, `waiting`, `requested`, or `pending` |
| conclusion | Final conclusion: `success`, `failure`, `neutral`, `cancelled`, `skipped`, `timed_out`, or `action_required` |
| event | Event that triggered the workflow run (for example, `push` or `pull_request`). Refer to [Events that trigger workflows](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows). |
| workflow_id | Unique identifier for the workflow definition |
| run_number | The run number for this workflow run in the repository |

### Workflow usage

Retrieve usage statistics for a workflow, such as run counts and durations.

#### Query options

| Name | Description | Required |
|------|-------------|----------|
| Owner | GitHub user or organization that owns the repository | Yes |
| Repository | Name of the repository | Yes |
| Workflow | The workflow ID or file name. Use `id` or the filename from `path` from [Workflows](#workflows) queries. | Yes |

##### Sample queries

Show usage statistics for the `Levitate` detect breaking changes workflow in the `grafana/grafana` repository:

- Owner: `grafana`
- Repository: `grafana`
- Workflow: `detect-breaking-changes-levitate.yml`

#### Response

| Name | Description |
|------|-------------|
| name | Name of the workflow (or workflow job) |
| unique triggering actors | Number of unique users or actors who triggered runs of this workflow |
| runs | Total number of workflow runs in the selected period |
| current billing cycle cost (approx.) | Approximate cost for the current billing cycle (if applicable) |
| skipped | Number (and percentage) of runs that were skipped |
| successes | Number (and percentage) of successful runs |
| failures | Number (and percentage) of failed runs |
| cancelled | Number (and percentage) of cancelled runs |
| total run duration (approx.) | Total duration of all runs (formatted as hours, minutes, seconds) |
| longest run duration (approx.) | Duration of the longest single run |
| average run duration (approx.) | Average duration of all runs |
| p95 run duration (approx.) | 95th percentile run duration |
| runs on Sunday | Number of runs started on Sunday |
| runs on Monday | Number of runs started on Monday |
| runs on Tuesday | Number of runs started on Tuesday |
| runs on Wednesday | Number of runs started on Wednesday |
| runs on Thursday | Number of runs started on Thursday |
| runs on Friday | Number of runs started on Friday |
| runs on Saturday | Number of runs started on Saturday |
