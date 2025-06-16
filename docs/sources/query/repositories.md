---
title: Repositories
menuTitle: Repositories
description: Repositories
hero:
  title: List github repositories
  level: 1
  width: 110
  image: https://raw.githubusercontent.com/grafana/github-datasource/refs/heads/main/src/img/github.svg
  height: 110
  description: List github repositories for a github organization / user
keywords:
  - github
  - ci
  - cd
labels:
  products:
    - oss
weight: 301
---

<!-- markdownlint-configure-file { "MD013": false, "MD033": false} -->

{{< docs/hero-simple key="hero" >}}

<hr style="margin-bottom:30px"/>

Query type **Repositories** allow you to list the repositories available for a github user/organization

## Query Parameters

| Field      | Description                                                                                                                                                               | Example                                                                |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| Query Type | Specify the query type as **Repository**                                                                                                                                  |                                                                        |
| Owner      | Github user id / organization id                                                                                                                                          | grafana                                                                |
| Repository | **(Optional)**. Search field. Typically full name of the repository or partial search string.<br/>**Tip**: You can use the same syntax you use in github.com search field | `datasource` / `grafana-infinity-datasource` / `datasources -lang:go`. |

### Query Schema

```json
{
  "queryType": "Repositories",
  "owner": "grafana",
  "repository": "datasource"
}
```

#### Downstream Query

Under the hood, grafana uses the following GraphQL query to list the repositories

```graphql
{
  search(query: "org:grafana datasource", type: REPOSITORY, first: 100) {
    nodes {
      ... on Repository {
        name
        owner {
          login
        }
        nameWithOwner
        url
        forks {
          totalCount
        }
        isFork
        isMirror
        isPrivate
        createdAt
      }
    }
  }
}
```

## Result field

Refer the [Github GraphQL documentation](https://docs.github.com/en/graphql/overview/explorer) for the details about all the result fields presented in the GraphQL query.

## Known Limitations

- Github plugin does pagination over the graphQL query with 100 pages limit on each request. If the owner/user have more number of repositories, then there are potential chances that the result might be slower. Use appropriate search filer in the repository field to increase the query performance.
