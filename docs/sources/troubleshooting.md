---
aliases:
  - ./troubleshoot/
title: Troubleshoot GitHub data source issues
menuTitle: Troubleshooting
description: Troubleshoot common issues with the GitHub data source plugin for Grafana
keywords:
  - data source
  - github
  - troubleshooting
  - errors
  - authentication
  - rate limiting
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 500
review_date: "2026-03-11"
---

# Troubleshoot GitHub data source issues

This page provides solutions to common issues you may encounter when configuring or using the GitHub data source. For configuration instructions, refer to [Configure the GitHub data source](https://grafana.com/docs/plugins/grafana-github-datasource/latest/configure/).

## Authentication errors

These errors occur when credentials are invalid, missing, or don't have the required permissions.

### "401 Unauthorized. Check your API key/Access token"

**Symptoms:**

- **Save & test** fails with a 401 error.
- Queries return authorization errors.
- Resources don't load in drop-downs.

**Possible causes and solutions:**

| Cause | Solution |
|-------|----------|
| Missing or incorrect token scopes | Verify the token includes all [required scopes](https://grafana.com/docs/plugins/grafana-github-datasource/latest/configure/#personal-access-token-permissions). |
| Expired or revoked token | Generate a new token in GitHub and update the data source configuration. |
| Extra whitespace or line breaks in token | Paste the token directly, avoiding extra spaces or line breaks. |
| Token doesn't have access to the queried repositories | Ensure the token has access to the repositories and organizations you're querying. For private repositories, the `repo` scope is required. |

### GitHub App authentication fails

**Symptoms:**

- **Save & test** fails when using GitHub App authentication.
- Queries return errors about invalid credentials.

**Possible causes and solutions:**

| Cause | Solution |
|-------|----------|
| Incorrect App ID | Verify the App ID matches the ID shown on your GitHub App's settings page. |
| Wrong Installation ID | Verify the installation ID. Navigate to **Settings** > **Installed GitHub Apps** > **Configure** and check the number at the end of the URL. |
| Expired or invalid private key | Generate a new private key from your GitHub App settings and update the data source configuration. |
| App not installed on the target organization | [Install the GitHub App](https://docs.github.com/en/apps/using-github-apps/installing-your-own-github-app) on the organization or user account you want to query. |
| Insufficient app permissions | Verify the app has the [required permissions](https://grafana.com/docs/plugins/grafana-github-datasource/latest/configure/#github-app-permissions). |

### Can't access private repositories or organizations

**Symptoms:**

- Queries only return data from public repositories.
- Organization-level queries return empty results.

**Solutions:**

1. Verify your personal access token includes the `repo` and `read:org` scopes.
1. For GitHub Apps, verify the app is installed on the organization and has the required repository permissions.
1. Generate a new token or update the app permissions and reconfigure the data source.

## Connection errors

These errors occur when Grafana can't reach the GitHub API endpoints.

### "404 Not Found. Check the Github Enterprise Server URL"

**Symptoms:**

- **Save & test** fails with a 404 error when using GitHub Enterprise Server.

**Solutions:**

1. Verify the GitHub Enterprise Server URL is correct and includes the protocol (for example, `https://github.example.com`).
1. Ensure the URL doesn't include a trailing path like `/api/v3` -- the plugin adds this automatically.
1. Verify the GitHub Enterprise Server instance is accessible from the Grafana server.

### "Unable to reach the Github Enterprise Server URL from the Grafana server"

**Symptoms:**

- **Save & test** fails with a connection error.
- DNS resolution or network timeout errors appear in Grafana logs.

**Solutions:**

1. Verify network connectivity from the Grafana server to the GitHub Enterprise Server.
1. Check that firewall rules allow outbound HTTPS (port 443) traffic to the GitHub endpoint.
1. Verify proxy settings if your Grafana server uses a proxy for outbound requests.
1. For Grafana Cloud, configure [Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/) if your GitHub Enterprise Server is in a private network.

## Query errors

These errors occur when executing queries against the data source.

### No data or empty results

**Symptoms:**

- Query executes without error but returns no data.
- Charts show "No data" message.

**Possible causes and solutions:**

| Cause | Solution |
|-------|----------|
| Time range doesn't contain data | Expand the dashboard time range. For queries without a time field (for example, time field set to `None`), the time range filter isn't applied. |
| Wrong owner or repository | Verify the **Owner** and **Repository** fields are correct. |
| Permissions issue | Verify the token or GitHub App has read access to the target repository. |
| Query syntax error | For Issues, Pull Requests, and Pull Request Reviews, verify the [GitHub search syntax](https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests). |

### Result limits

Some query types have maximum result limits:

| Query type | Maximum results | Source |
|-----------|----------------|--------|
| Issues | 1,000 | GitHub search API limit |
| Repositories | 1,000 | GitHub search API limit |
| Contributors | 200 | Plugin page limit |
| Projects | 200 | Plugin page limit |

If you're hitting these limits, use more specific query filters or narrow the time range to reduce the result set.

### Rate limiting

**Symptoms:**

- Queries intermittently fail or return errors.
- Multiple panels on a dashboard fail to load simultaneously.

**Solutions:**

The plugin caches all API responses for up to five minutes to reduce the number of requests to GitHub. If you're still hitting rate limits:

1. Reduce the frequency of dashboard auto-refresh.
1. Use more specific queries to reduce the number of API calls.
1. If using a personal access token, consider switching to a GitHub App, which has higher rate limits.
1. For Grafana Enterprise or Grafana Cloud, enable [query caching](https://grafana.com/docs/grafana/latest/administration/data-source-management/#query-and-resource-caching) for additional caching control beyond the built-in cache.

## Code scanning errors

### Code scanning alerts return empty results

**Symptoms:**

- Code scanning queries return no data even though alerts exist in GitHub.

**Solutions:**

1. Verify the token or GitHub App has the [code scanning permissions](https://grafana.com/docs/plugins/grafana-github-datasource/latest/configure/#code-scanning-permissions): `code scanning alerts: read-only` and `security_events: read-only`.
1. Verify code scanning is enabled on the repository.
1. If querying at the organization level, leave the **Repository** field empty.

## Enable debug logging

To capture detailed error information for troubleshooting:

1. Set the Grafana log level to `debug` in the configuration file:

   ```ini
   [log]
   level = debug
   ```

1. Restart Grafana for the change to take effect.
1. Reproduce the issue and review logs in `/var/log/grafana/grafana.log` (or your configured log location).
1. Look for entries containing `github` for request and response details.
1. Reset the log level to `info` after troubleshooting to avoid excessive log volume.

## Get additional help

If you've tried the solutions on this page and still have issues:

1. Check the [Grafana community forums](https://community.grafana.com/) for similar issues.
1. Review the [GitHub data source plugin issues](https://github.com/grafana/github-datasource/issues) for known bugs.
1. Consult the [GitHub API documentation](https://docs.github.com/en/rest) for service-specific guidance.
1. Contact [Grafana Support](https://grafana.com/support/) if you're a Grafana Cloud Pro, Cloud Advanced, or Enterprise customer.
1. When reporting issues, include:
   - Grafana version and plugin version.
   - Exact error messages (redact sensitive information).
   - Steps to reproduce the issue.
   - Relevant configuration details (redact credentials).
