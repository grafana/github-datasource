---
title: Use macros with GitHub data source plugin for Grafana
menuTitle: Macros
description: Learn about the macros you can use in the GitHub data source plugin for Grafana
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
weight: 302
---

# Use macros with GitHub data source plugin for Grafana

A macro is a feature that allows you to simplify the syntax and add dynamic parts to your queries.
They help make your queries more flexible.

The GitHub data source plugin for Grafana supports the following macros:

| Macro name | Syntax                     | Description                                                              | Example                                                                              |
| ---------- | -------------------------- | ------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| multiVar   | `$__multiVar(prefix,$var)` | Expands a multi value variable into github query string                  | `$__multiVar(label,$labels)` will expand into `label:first-label label:second-label` |
|            |                            | When using **all** in multi variable, use **\*** as custom all value     |                                                                                      |
| day        | `$__toDay(diff)`           | Returns the day according to UTC time, a difference in days can be added | `created:$__toDay(-7)` on 2022-01-17 will expand into `created:2022-01-10`           |
