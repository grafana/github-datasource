# Grafana GitHub datasource

⚠️ Work in Progress⚠️

The GitHub datasource allows GitHub API data to be visually represented in Grafana dashboards.

## Features

### Backend
* [x] Releases
  * [x] List all
  * [x] List within time range
  * [ ] Filter by author
  * [ ] Filter by semver
* [x] Commits
  * [x] List all
  * [x] List within time range
  * [ ] Filter by author
* [ ] Deploys
* [x] Repositories
  * [x] List all
* [x] Issues
  * [x] List all
  * [x] List within time range
  * [x] Filter
    * [x] By Assignee
    * [x] Opened at
    * [x] Closed at
    * [x] By labels
* [x] Organizations
* [ ] Milestones
* [ ] Response Caching

### Frontend
* [ ] Visualize queries
* [ ] Template variables
* [ ] Alerts / Notification Channels
* [ ] Annotations
* [ ] Authorize requests
  * [ ] OAuth
  * [ ] Personal Access Token
