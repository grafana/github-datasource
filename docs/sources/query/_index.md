# GitHub Data Source Plugin Overview

The GitHub data source plugin for Grafana enables you to query and visualize data directly from your GitHub repositories and organizations. With this plugin, you can monitor repository activity, track issues and pull requests, analyze workflow runs, and more—all within your Grafana dashboards.

## Query Types

The plugin supports the following query types:

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
- **Branch**: (Optional) The branch to filter commits by.
- **Author**: (Optional) Filter commits by a specific author.
- **Since/Until**: (Optional) Filter commits within a date range.

**Sample query:**  
Show all commits to the `main` branch of the `grafana/grafana` repository in the last 7 days.

- Owner: `grafana`
- Repository: `grafana`
- Branch: `main`
- Since: `2024-05-08`
- Until: `2024-05-15`

**Sample response:**
```json
[
  {
    "sha": "7638417db6d59f3c431d3e1f261cc637155684cd",
    "commit": {
      "author": {
        "name": "Monalisa Octocat",
        "email": "octocat@github.com",
        "date": "2024-05-10T12:34:56Z"
      },
      "message": "Fix all the bugs"
    },
    "author": {
      "login": "octocat",
      "id": 1
    },
    "html_url": "https://github.com/grafana/grafana/commit/7638417db6d59f3c431d3e1f261cc637155684cd"
  }
]
```

---

### Issues

List issues in a repository, with options to filter by state, assignee, and labels. Useful for tracking open bugs, feature requests, or project tasks.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.
- **State**: Filter by issue state (`open`, `closed`, or `all`).
- **Assignee**: (Optional) Filter by assigned user.
- **Labels**: (Optional) Filter by one or more labels.
- **Creator**: (Optional) Filter by issue creator.
- **Since**: (Optional) Only issues updated at or after this time.

**Sample query:**  
Show all open issues labeled `bug` in the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`
- State: `open`
- Labels: `bug`

**Sample response:**
```json
[
  {
    "id": 1,
    "number": 1347,
    "title": "Found a bug",
    "user": {
      "login": "octocat"
    },
    "state": "open",
    "labels": [
      {
        "name": "bug"
      }
    ],
    "created_at": "2024-05-12T09:00:00Z",
    "html_url": "https://github.com/grafana/grafana/issues/1347"
  }
]
```

---

### Contributors

Get a list of contributors to a repository, including their contribution counts. Useful for understanding project activity and top contributors.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all contributors to the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

**Sample response:**
```json
[
  {
    "login": "octocat",
    "id": 1,
    "contributions": 42
  }
]
```

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

**Sample response:**
```json
[
  {
    "name": "v10.0.0",
    "commit": {
      "sha": "abc123"
    }
  }
]
```

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

**Sample response:**
```json
[
  {
    "id": 1,
    "tag_name": "v10.0.0",
    "name": "First Release",
    "draft": false,
    "prerelease": false,
    "created_at": "2024-05-01T00:00:00Z",
    "html_url": "https://github.com/grafana/grafana/releases/tag/v10.0.0"
  }
]
```

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

**Sample response:**
```json
[
  {
    "id": 2,
    "number": 42,
    "title": "Add new feature",
    "user": {
      "login": "octocat"
    },
    "state": "open",
    "created_at": "2024-05-13T10:00:00Z",
    "html_url": "https://github.com/grafana/grafana/pull/42"
  }
]
```

---

### Labels

Get all labels defined in a repository, useful for categorizing issues and pull requests.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all labels for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

**Sample response:**
```json
[
  {
    "id": 208045946,
    "name": "bug",
    "color": "f29513"
  }
]
```

---

### Repositories

List repositories for a user or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.

**Sample query:**  
Show all repositories for the `grafana` organization.

- Organization: `grafana`

**Sample response:**
```json
[
  {
    "id": 1296269,
    "name": "grafana",
    "full_name": "grafana/grafana",
    "private": false,
    "html_url": "https://github.com/grafana/grafana"
  }
]
```

---

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

**Sample response:**
```json
[
  {
    "id": 1002604,
    "number": 1,
    "title": "v10.0.0",
    "state": "open",
    "description": "Release milestone for v10.0.0"
  }
]
```

---

### Packages

List packages published in a repository or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.
- **Repository**: (Optional) The name of the repository.

**Sample query:**  
Show all packages for the `grafana` organization.

- Organization: `grafana`

**Sample response:**
```json
[
  {
    "id": 1,
    "name": "grafana-package",
    "package_type": "container",
    "html_url": "https://github.com/orgs/grafana/packages/container/grafana-package"
  }
]
```

---

### Vulnerabilities

Query security vulnerabilities detected in a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all security advisories for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

**Sample response:**
```json
[
  {
    "ghsa_id": "GHSA-xxxx-yyyy-zzzz",
    "summary": "Vulnerability in dependency",
    "severity": "high",
    "url": "https://github.com/advisories/GHSA-xxxx-yyyy-zzzz"
  }
]
```

---

### Projects

List classic projects associated with a repository or organization.

#### Query options

- **Owner/Organization**: The GitHub user or organization.
- **Repository**: (Optional) The name of the repository.

**Sample query:**  
Show all projects for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

**Sample response:**
```json
[
  {
    "id": 1002604,
    "name": "Roadmap",
    "body": "Project roadmap for 2024",
    "state": "open"
  }
]
```

---

### Stargazers

Get a list of users who have starred a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all stargazers for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

**Sample response:**
```json
[
  {
    "starred_at": "2024-05-14T12:00:00Z",
    "user": {
      "login": "octocat",
      "id": 1
    }
  }
]
```

---

### Workflows

List GitHub Actions workflows defined in a repository.

#### Query options

- **Owner**: The GitHub user or organization that owns the repository.
- **Repository**: The name of the repository.

**Sample query:**  
Show all workflows for the `grafana/grafana` repository.

- Owner: `grafana`
- Repository: `grafana`

**Sample response:**
```json
{
  "total_count": 2,
  "workflows": [
    {
      "id": 161335,
      "name": "CI",
      "path": ".github/workflows/ci.yml",
      "state": "active"
    }
  ]
}
```

---

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

**Sample response:**
```json
{
  "billable": {
    "UBUNTU": {
      "total_ms": 1800000,
      "jobs": 3,
      "runs": 3
    }
  }
}
```

---

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

**Sample response:**
```json
{
  "total_count": 1,
  "workflow_runs": [
    {
      "id": 30433642,
      "name": "CI",
      "head_branch": "main",
      "status": "completed",
      "conclusion": "success",
      "run_number": 562,
      "created_at": "2024-05-14T10:00:00Z",
      "html_url": "https://github.com/grafana/grafana/actions/runs/30433642"
    }
  ]
}
```
---