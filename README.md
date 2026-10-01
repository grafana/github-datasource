# RWE Grafana GitHub data source

The GitHub data source plugin for Grafana lets you to query the GitHub API in Grafana so you can visualize your GitHub repositories and projects.

## Documentation

This fork adds GitHub Enterprise Cloud data-residency support to the upstream Grafana GitHub data source. See the [configuration guide](docs/sources/configure.md) for supported connection modes.

## Build and sign

Use the Node version from `.nvmrc` and the Go version declared in `go.mod`:

```bash
npm ci
npm run test:ci
npm run build
go run github.com/magefile/mage -v build:linux
```

Grafana requires a signature for production plugins. Set `GRAFANA_ACCESS_POLICY_TOKEN` in the release environment, then sign for the exact externally visible Grafana URLs:

```bash
npm run sign -- --rootUrls https://grafana.example.com/
```

Do not commit the signing token. Unsigned loading with `GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=rwe-github-datasource` is intended only for local development and test environments.

Follow [Release and deploy the private GitHub data source](docs/private-plugin-release.md) for the one-time Grafana/GitHub setup, release procedure, sphere rollout, datasource migration, verification, and token rotation steps.

## Video Tutorial

Watch this video to learn more about setting up the Grafana GitHub data source plugin:

[![GitHub data source plugin | Visualize GitHub using Grafana | Tutorial](https://img.youtube.com/vi/DW693S3cO48/hq720.jpg)](https://youtu.be/DW693S3cO48 "Grafana GitHub data source plugin")

> [!TIP]
> 
> ## Give it a try using Grafana Play
> 
> With Grafana Play, you can explore and see how it works, learning from practical examples to accelerate your development. This feature can be seen on [GitHub data source plugin demo](https://play.grafana.org/d/d5b56357-1a57-4821-ab27-16fdf79cab57/github3a-queries-and-multi-variables).

## GitHub API V4 (GraphQL)

This data source uses the [`githubv4` package](https://github.com/shurcooL/githubv4), which is under active development.

## Private data source connect - Only for Grafana Cloud users.

Establishes a private, secured connection between a Grafana Cloud stack and data sources within a private network. Use the drop-down to locate the PDC URL. For setup instructions, refer to [Private data source connect (PDC)](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/) and [Configure PDC](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/configure-pdc/#configure-grafana-private-data-source-connect-pdc). Click Manage private data source connect to open your PDC connection page and view your configuration details.

## Frequently Asked Questions

- **Why does it sometimes take up to 5 minutes for my new pull request / new issue / new commit to show up?**

We have aggressive caching enabled due to GitHub's rate limiting policies. When selecting a time range like "Last hour", a combination of the queries for each panel and the time range is cached temporarily.

- **Why are there two selection options for Pull Requests and Issue times when creating annotations?**

There are two times that affect an annotation:

- The time range of the dashboard or panel
- The time that should be used to display the event on the graph

The first selection is used to filter the events that display on the graph. For example, if you select "closed at", only events that were "closed" in your dashboard's time range will be displayed on the graph.

The second selection is used to determine where on the graph the event should be displayed.

Typically, these will be the same, however there are some cases where you may want them to be different.
