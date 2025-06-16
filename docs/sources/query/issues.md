---
title: Issues
menuTitle: Issues
description: Issues
hero:
  title: List issues
  level: 1
  width: 110
  image: https://raw.githubusercontent.com/grafana/github-datasource/refs/heads/main/src/img/github.svg
  height: 110
  description: List issues created in a github repo
keywords:
  - github
  - ci
  - cd
labels:
  products:
    - oss
weight: 302
---

<!-- markdownlint-configure-file { "MD013": false, "MD033": false} -->

{{< docs/hero-simple key="hero" >}}

<hr style="margin-bottom:30px"/>

Query type **Issues** allows you to list the issues created in a GitHub repository.

## Query Parameters

| Field      | Description                                                             | Example   |
| ---------- | ----------------------------------------------------------------------- | --------- |
| Query Type | Specify the query type as **Issues**                                    |           |
| Owner      | GitHub user id / organization id                                        | grafana   |
| Repository | Name of the repository                                                  | grafana   |
| Options    | **(Optional)**. Additional options such as search query and time field. | See below |

### Options

| Field          | Description                                                                                                                                                                                                                                                                   | Example                  |
| -------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| **Query**      | _(Optional)_ Search string to filter issues. You can use the same syntax as GitHub's search field.                                                                                                                                                                            | `is:OPEN label:type/bug` |
| **Time Field** | _(Optional)_ Field to use for time filtering (`created`=0, `closed`=1, `updated`=2). If nothing selected, defaults to `created at` and the results will be filtered automatically based on the date. <br/>**Tip**: Grafana's dashboard time range used for these time filters | 0                        |

### Query Schema

```json
{
  "queryType": "Issues",
  "owner": "grafana",
  "repository": "grafana",
  "options": {
    "query": "is:open label:bug",
    "timeField": "created"
  }
}
```

#### Downstream Query

Under the hood, Grafana uses the following GraphQL query to list the issues:

```graphql
TBD
```

## Result field

Refer to the [GitHub GraphQL documentation](https://docs.github.com/en/graphql/overview/explorer) for details about all the result fields presented in the GraphQL query.

## Known Limitations

- The GitHub plugin paginates over the GraphQL query with a 100 items limit per request. If the repository has a large number of issues, the result might be slower. Use appropriate filters (query, timeField) to improve query performance.
