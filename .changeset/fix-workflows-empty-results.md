---
"grafana-github-datasource": patch
---

Fix empty results in Workflows query type. Added nil check for CreatedAt/UpdatedAt timestamps and added "None" option to Time Field dropdown (default) to return all workflows without time filtering.
