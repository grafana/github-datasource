---
'grafana-github-datasource': patch
---

Report GitHub App installation token refresh failures and GitHub Enterprise token policy rejections as downstream errors, so a data source with invalid credentials no longer counts against the plugin error budget
