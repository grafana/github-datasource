# Grafana GitHub data source

The GitHub data source plugin for Grafana lets you to query the GitHub API in Grafana so you can visualize your GitHub repositories and projects.

## Documentation

For the plugin documentation, visit <a href="https://grafana.com/docs/plugins/grafana-github-datasource" target="_blank">plugin documentation website</a>

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
