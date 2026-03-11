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

Use template variables to create dynamic, reusable dashboards with the GitHub data source.

A [variable](https://grafana.com/docs/grafana/latest/dashboards/variables/) is a placeholder for a value that you can use in dashboard queries. Variables are displayed as drop-down lists at the top of the dashboard, making it easy to change the data being displayed without editing queries directly.

## Create a query variable

To create a query variable using the GitHub data source:

1. Navigate to **Dashboard settings** > **Variables**.
1. Click **Add variable**.
1. Select **Query** as the variable type.
1. Select the GitHub data source from the **Data source** drop-down.
1. Select the **Query Type** (for example, Repositories, Labels, or any other supported query type).
1. Configure the query options for the selected query type.
1. Use **Field Value** to select which response field provides the variable values.
1. Use **Field Display** to select which response field provides the display labels.

All query types supported by the GitHub data source can be used as variable queries.

## Use variables in queries

You can reference variables in any string field in the query editor, such as **Owner**, **Repository**, or **Query**, using the standard Grafana variable syntax:

- `$variable` or `${variable}` for single-value variables.
- Variables with multiple selected values are interpolated as comma-separated values.

For example, to create a dashboard that lets users select a repository:

1. Create a variable named `repo` using the **Repositories** query type.
1. In your query, set **Repository** to `$repo`.

## Macros

The GitHub data source provides macros that add dynamic parts to your queries. Use macros in the **Query** field.

| Macro | Syntax | Description | Example |
|-------|--------|-------------|---------|
| **multiVar** | `$__multiVar(prefix,$var)` | Expands a multi-value variable into GitHub query qualifiers. Each selected value is prefixed with the specified prefix. | `$__multiVar(label,$labels)` expands to `label:first-label label:second-label` |
| **toDay** | `$__toDay(diff)` | Returns the current date in UTC. An optional `diff` parameter offsets the date by the specified number of days. | `created:$__toDay(-7)` on 2022-01-17 expands to `created:2022-01-10` |

### Use multiVar with "All" selection

When using the `$__multiVar` macro with a multi-value variable, set the **Custom all value** to `*` in the variable settings. When "All" is selected, the macro expands to an empty string, which effectively removes the filter from the query.

### Macro examples

Filter issues by multiple labels selected from a variable:

- Query: `is:open $__multiVar(label,$labels)`

Filter issues created in the last 7 days:

- Query: `is:open created:>=$__toDay(-7)`

Combine macros for a dynamic query:

- Query: `is:open $__multiVar(label,$labels) created:>=$__toDay(-30)`
