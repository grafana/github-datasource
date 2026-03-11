---
aliases:
  - ./variables-and-macros/
  - ./variables-and-macros/variables/
  - ./variables-and-macros/macros/
title: GitHub template variables
menuTitle: Template variables
description: Use template variables and macros with the GitHub data source plugin for Grafana
keywords:
  - data source
  - github
  - variables
  - macros
  - template
  - dynamic dashboards
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 300
review_date: "2026-03-11"
---

# GitHub template variables

Use template variables to create dynamic, reusable dashboards with the GitHub data source. For general information about variables in Grafana, refer to [Variables](https://grafana.com/docs/grafana/latest/dashboards/variables/).

## Before you begin

- [Configure the GitHub data source](https://grafana.com/docs/plugins/grafana-github-datasource/latest/configure/).
- Familiarize yourself with the available [query types](https://grafana.com/docs/plugins/grafana-github-datasource/latest/query-editor/).

## Supported variable types

The GitHub data source supports **Query** variables. You can use any query type (for example, Repositories, Labels, Issues) to populate variable options.

## Create a query variable

To create a query variable:

1. Navigate to **Dashboard settings** > **Variables**.
1. Click **Add variable**.
1. Select **Query** as the variable type.
1. Select the GitHub data source from the **Data source** drop-down.
1. Select the **Query Type** (for example, Repositories, Labels, or any other supported query type).
1. Configure the query options for the selected query type.
1. Set **Field Value** to the response field that provides the variable values (the actual value stored when a user selects an option).
1. Set **Field Display** to the response field that provides the display labels (the text shown in the drop-down).

## Use variables in queries

You can reference variables in any string field in the query editor, such as **Owner**, **Repository**, or **Query**, using the standard Grafana variable syntax:

- `$variable` or `${variable}` for single-value variables.
- Variables with multiple selected values are interpolated as comma-separated values.

For example, to create a dashboard that lets users select a repository:

1. Create a variable named `repo` using the **Repositories** query type.
1. In your query, set **Repository** to `$repo`.

## Macros

The GitHub data source provides macros that add dynamic parts to your queries. Macros are used in the **Query** field, which is available on query types that support [GitHub search syntax](https://docs.github.com/en/search-github/getting-started-with-searching-on-github/understanding-the-search-syntax), such as Issues, Pull Requests, and Pull Request Reviews.

| Macro | Syntax | Description | Example |
|-------|--------|-------------|---------|
| **multiVar** | `$__multiVar(qualifier,$var)` | Expands a multi-value variable into GitHub search qualifiers. Each selected value is combined with the qualifier name in `qualifier:value` format. | `$__multiVar(label,$labels)` expands to `label:first-label label:second-label` |
| **toDay** | `$__toDay(diff)` | Returns the current date in UTC, formatted as `YYYY-MM-DD`. An optional `diff` parameter offsets the date by the specified number of days. | `created:$__toDay(-7)` on 2022-01-17 expands to `created:2022-01-10` |

### Use multiVar with "All" selection

When using the `$__multiVar` macro with a multi-value variable, set the **Custom all value** to `*` in the variable settings. When "All" is selected, the macro returns an empty string, which removes that qualifier from the query.

### Macro examples

Filter issues by multiple labels selected from a variable:

- Query: `is:open $__multiVar(label,$labels)`

Filter issues created in the last 7 days:

- Query: `is:open created:>=$__toDay(-7)`

Combine macros for a dynamic query:

- Query: `is:open $__multiVar(label,$labels) created:>=$__toDay(-30)`

## Use case: Multi-variable dashboard

You can chain variables together to build a dashboard that lets users drill down from an organization to a specific repository and filter by labels.

### Set up the variables

Create the following variables in order. Later variables reference earlier ones, so the order matters.

1. **`owner`** — Type: **Custom**. Add a comma-separated list of your GitHub organizations or user accounts (for example, `grafana,prometheus`). This provides the top-level organization selector.

1. **`repo`** — Type: **Query**. Use the **Repositories** query type with **Owner** set to `$owner`. Set **Field Value** and **Field Display** to `name`. This populates the repository drop-down based on the selected organization.

1. **`labels`** — Type: **Query**, **Multi-value** enabled. Use the **Labels** query type with **Owner** set to `$owner` and **Repository** set to `$repo`. Set **Field Value** and **Field Display** to `name`. Set **Custom all value** to `*`. This lets users filter by one or more labels.

### Use the variables in panels

With these variables in place, build panels that respond to all three selectors:

- **Open issues panel:** Set **Owner** to `$owner`, **Repository** to `$repo`, and **Query** to `is:open $__multiVar(label,$labels)`.
- **Pull request panel:** Set **Owner** to `$owner`, **Repository** to `$repo`, and **Query** to `is:open $__multiVar(label,$labels)`.
- **Commits panel:** Set **Owner** to `$owner` and **Repository** to `$repo` to show commit activity for the selected repository.

When a user changes the organization drop-down, the repository list refreshes automatically. Selecting a different repository updates the label list and all panels on the dashboard.
