---
slug: '/macros'
title: 'Macros'
menuTitle: Macros
description: Using macros for the GitHub data source plugin
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
weight: 202
---

# Macros

A Macro is a feature that allows you to simplify the syntax and add dynamic parts to your queries. They help make your queries more flexible and adaptable to changing conditions. It's important to note that while not all data source plugins have macros as it availability and implementation depends on the specific design and capabilities of each plugin.

The GitHub data source plugin do support the macro feature and you can use the following macros in your queries:

| Macro name | Syntax                     | Description                                                              | Example                                                                              |
| ---------- | -------------------------- | ------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| multiVar   | `$__multiVar(prefix,$var)` | Expands a multi value variable into github query string                  | `$__multiVar(label,$labels)` will expand into `label:first-label label:second-label` |
|            |                            | When using **all** in multi variable, use **\*** as custom all value     |                                                                                      |
| day        | `$__toDay(diff)`           | Returns the day according to UTC time, a difference in days can be added | `created:$__toDay(-7)` on 2022-01-17 will expand into `created:2022-01-10`           |