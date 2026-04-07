---
title: Troubleshooting common problems for the Grafana GitHub data source plugin
menuTitle: Troubleshoot
description: Learn how to troubleshoot common problems for the GitHub data source plugin
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
weight: 600
---

# Troubleshooting common problems for the Grafana GitHub data source plugin

This page lists some common issues you may experience when setting up the Grafana GitHub data source plugin. You can check the possible reasons and suggested solutions below.

### Why does my data source setup fail ("Authentication failed") when using a GitHub Personal Access Token (PAT)?

- Make sure your PAT is created with valid [permissions required for the Grafana GitHub data source](https://grafana.com/docs/plugins/grafana-github-datasource/setup/token/#permissions).
- Check that the token is not expired or revoked.
- Ensure the token has access to the repositories you want to query
- Paste the PAT directly into the configuration field, avoiding extra spaces or line breaks.
- If using GitHub Enterprise, verify the API URL and ensure the PAT is valid for that instance.
- If the error persists, generate a new PAT and update the configuration.

### Why can't I access private repositories or organizations?

- Make sure your PAT includes the `repo` and `read:org` scopes.
- Update the data source configuration with the new token and test again.

### What should I do if I see "An unexpected error happened" or "Could not connect to GitHub" after trying all of the above?

- Check the Grafana logs for more details about the error.
- For Grafana Cloud customers, contact support.
