## Variables

[Variables](https://grafana.com/docs/grafana/latest/variables/) allow you to substitute values in a panel with pre-defined values.

![Creating Variables](https://github.com/grafana/github-datasource/raw/main/docs/screenshots/variables-create.png)

You can reference them inside queries, allowing users to configure parameters such as `Query` or `Repository`.

![Using Variables inside queries](https://github.com/grafana/github-datasource/raw/main/docs/screenshots/using-variables.png)

## Macros

You can use the following macros in your queries

| Macro Name | Syntax                     | Description                                                              | Example                                                                              |
| ---------- | -------------------------- | ------------------------------------------------------------------------ | ------------------------------------------------------------------------------------ |
| multiVar   | `$__multiVar(prefix,$var)` | Expands a multi value variable into github query string                  | `$__multiVar(label,$labels)` will expand into `label:first-label label:second-label` |
|            |                            | When using **all** in multi variable, use **\*** as custom all value     |                                                                                      |
| day        | `$__toDay(diff)`           | Returns the day according to UTC time, a difference in days can be added | `created:$__toDay(-7)` on 2022-01-17 will expand into `created:2022-01-10`           |