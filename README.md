# Grafana GitHub data source

The GitHub data source plugin for Grafana lets you to query the GitHub API in Grafana so you can visualize your GitHub repositories and projects.

## Documentation

For the plugin documentation, visit [plugin documentation website](https://grafana.com/docs/plugins/grafana-github-datasource).

## GitHub API V4 (GraphQL)

This data source uses the [`githubv4` package](https://github.com/shurcooL/githubv4), which is under active development.

## Frequently Asked Questions

- **I am using GitHub OAuth on Grafana. Can my users make requests with their individual GitHub accounts instead of a shared `access_token`?**

No. This requires changes in Grafana first. See [this issue](https://github.com/grafana/grafana/issues/26023) in the Grafana project.

- **Why does it sometimes take up to 5 minutes for my new pull request / new issue / new commit to show up?**

We have aggressive caching enabled due to GitHub's rate limiting policies. When selecting a time range like "Last hour", a combination of the queries for each panel and the time range is cached temporarily.

- **Why are there two selection options for Pull Requests and Issue times when creating annotations?**

There are two times that affect an annotation:

- The time range of the dashboard or panel
- The time that should be used to display the event on the graph

The first selection is used to filter the events that display on the graph. For example, if you select "closed at", only events that were "closed" in your dashboard's time range will be displayed on the graph.

The second selection is used to determine where on the graph the event should be displayed.

Typically these will be the same, however there are some cases where you may want them to be different.
