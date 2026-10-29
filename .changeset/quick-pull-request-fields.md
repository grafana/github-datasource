---
'grafana-github-datasource': minor
---

adds an Additional fields option to the Pull Requests query; clearing the fields that are slow or large for GitHub to return (additions / deletions, mergeable, author, merged by, labels) makes queries on large repositories several times faster. All fields are still fetched by default
