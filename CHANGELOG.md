# Change Log

## 2.3.0

🚀 Add UpdatedAt time field to pull request queries

🐛 Add runStartedAt field for workflow runs query response

## 2.2.0

🚀 Adds support for Pull Request Review queries

⚙️ Bump form-data from 4.0.0 to 4.0.4 (#490)

⚙️ Add alternate "rate limit exceeded" error handling (#494)

🐛 Remove <base target="_blank"> from readme (#497)

⚙️ Bump @grafana/create-plugin configuration to 5.25.8 (#496)

## 2.1.7

🐛 Return assignees with issues queries

## 2.1.6

🐛 Documentation links will open in a new tab
🐛 Removed unused annotations method (replaced with new annotations support in [#196](https://github.com/grafana/github-datasource/pull/196))
🐛 Fixes error parsing app id / client id through provisioning via environment variables. Fixes [#477](https://github.com/grafana/github-datasource/issues/477)
🐛 Replaced the deprecated `setVariableQueryEditor` with `CustomVariableSupport`

## 2.1.5

🐛 Update golang-jwt/jwt dependency to v4.5.2
🐛 Get default http transport from plugin-sdk-go

## [2.1.4]

- **Fix** - Workflow runs - date filter now filters by time
- **Fix** - Panic in project query
- **Security** - Bump prismjs to 1.30.0

## [2.1.3]

- **Fix** - Add mutex protection to prevent data races in datasource cache
- **Chore** - Add validation for package types

## [2.1.2]

- **Fix** - GitHub enterprise url missing /api/v3

## [2.1.1]

- **Fix** - GitHub enterprise url wrong with app authentication
- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go to 0.268.1

## [2.1.0]

- **Feature** - Add a new query to retrieve Workflow runs
- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go to 0.266.0
- **Chore** - Bump dompurify to 3.2.4

## [2.0.2]

- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go to 0.265.0

## [2.0.1]

- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go from 0.260.3 to 0.261.0

## [2.0.0]

- **Chore**: Plugin now requires Grafana 10.4.8 or newer

## [1.9.2]

- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go from 0.255.0 to 0.258.0
- **Chore** - Update uplot dependency to 1.6.31

## [1.9.1]

- **Docs** - Add video tutorial to the README
- **Docs** - Update documentation for permissions and provisioning
- **Docs** - Update documentation plugin setup
- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go from 0.252.0 to 0.255.0
- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go from 0.251.0 to 0.252.0
- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go from 0.250.0 to 0.251.0

## [1.9.0]

- **Feature** - Add support for GitHub App authentication
- **Fix** - Fix error source for saml error and limit error
- **Fix** - Hide package types that are not supported by the GraphQL API
- **Chore** - Update spelling of GitHub
- **Chore** - Bump github.com/grafana/grafana-plugin-sdk-go from 0.247.0 to 0.250.0
- **Chore** - Bump path-to-regexp from 1.8.0 to 1.9.0

## [1.8.2]

- **Chore** - Bump grafana-plugin-sdk-go to v0.247.0

## [1.8.1]

- **Chore** - Bump micromatch from 4.0.5 to 4.0.8
- **Chore** - Bump webpack from 5.91.0 to 5.94.0
- **Chore** - Bump grafana-plugin-sdk-go to v0.245.0

## [1.8.0]

- **Feature** - Add additional user field to results in `Pull Request` query
- **Chore** - Update documentation, new and updated documentation available on official website

## [1.7.4]

- **Chore** - Update documentation

## [1.7.3]

- **Fix** - Fix error in `Packages` query where no package type was initially selected
- **Chore** - Update documentation
- **Chore** - Bump grafana-plugin-sdk-go to v0.241.0
- **Chore** - Mark downstream errors

## [1.7.2]

- **Chore** - Bump grafana-plugin-sdk-go to v0.240.0

## [1.7.1]

- **Chore** - Add provisioning folder to .gitignore
- **Chore** - Add error source to error response

## [1.7.0]

- **Feature** - Add `updated_at` field to results in `Issue` query
- **Feature** - Add `UpdatedAt` field to query options in `Issue` query
- **Fix** - Fix error when response has data with empty array in templating
- **Fix** - Fix per page limit to 100 in `Workflows` query as it is max supported value
- **Fix** - Remove query input in `Vulnerabilities` query as API does not support it
- **Chore** - Move e2e from cypress to playwright
- **Chore** - Update dependencies

## [1.6.0]

- **Feature** - Add `message` field to `Commit` query
- **Feature** - Add `name` field to `Workflow status` query
- **Fix** - Variable editor to support all query types

## [1.5.7]

- **Chore** - Update dependencies

## [1.5.6]

- **Chore** - Build with go 1.22.2
- **Chore** - Bump grafana-plugin-sdk-go to v0.220.0 (latest)
- **Bug Fix** - Prevent partial queries running on change of query type

## [1.5.5]

- **Chore** - Build with go 1.22
- **Fix** - Make health check faster by using github-datasource repository instead of grafana

## [1.5.4]

- **Chore** - Bump grafana-plugin-sdk-go to v0.198.0 (latest)
- **Bug Fix** - Fix tag queries to return commits as well
- **Bug Fix** - Fix for resetting URL in the config page

## [1.5.3]

- **Chore** - Bump grafana-plugin-sdk-go to latest
- **Chore** - Added lint GitHub workflow
- **Chore** - Remove legacy form styling

## [1.5.2]

- **BugFix** - Fix config page backwards compatibility with Grafana < 10.1

## [1.5.1] - 2023-10-10

- **Feature** - Update configuration page
- **Chore** - Update feature tracking usage to improve performance

## [1.5.0] - 2023-09-13

- **Feature** - Issues Query: Allow repo to be optional

## [1.4.7] - 2023-08-03

- **Feature** - Add ability to query Workflow and Workflow usage

## [1.4.6] - 2023-07-14

- **Bugfix** - Fixed a bug where disabled queries were still being executed

## [1.4.5] - 2023-05-04

- **Chore** - Backend binaries are now compiled with golang 1.20.4

## [1.4.4] - 2023-04-19

- **Chore** - Updated go version to 1.20

## [1.4.3] - 2023-03-07

- **Chore** - Update grafana-plugin-sdk-go to v0.155.0 to fix `The content of this plugin does not match its signature` error

## [1.4.2] - 2023-03-06

- **Chore** - Migrate to create plugin and upgrade dependencies

## [1.4.1] - 2023-03-01

- **Feature** - Added `RepositoryVulnerabilityAlertState` field to `Vulnerabilities` query

## [1.4.0] - 2023-02-03

- **Feature** - Added stargazers query type
- **Chore** - Minor documentation updates

## [1.3.3] - 2023-01-09

- **Chore** - Removed angular dependency: migrated annotation editor

## [1.3.2] - next

- **Feature** Added `$__toDay()` macro support

## [1.3.1] 2022-12-21

- **Chore** - Updated go version to latest (1.19.4)
- **Chore** - Updated backend grafana dependencies
- **Chore** - Added spellcheck

## [1.3.0] 2022-11-3

- **Feature** - GitHub projects - query items, user projects
- **Chore** - Updated build to use go 1.19.3

## [1.2.0] 2022-10-20

- **Feature** - GitHub projects

## [1.1.0] - next

- Updated grafana minimum runtime required to 8.4.7

## [1.0.15] 2022-05-05

- Fix variable interpolation

## [1.0.14] 2022-04-25

- Added a `$__multiVar()` macro support

## [1.0.13] 2021-12-01

- Fixed a bug where dashboard variables could not be set properly

## [1.0.12] 2021-12-01

- Added refId in annotation queries

## [1.0.11] 2021-05-17

- Added repository fields to the responses

## [1.0.10] 2021-04-01

- Fixed issue where some time values were being rendered incorrectly

## [1.0.9] 2021-04-01

- Fixed issue where dashboard path was not incorrect

## [1.0.8] 2020-12-10

- Fixed issue where screenshots were not rendering on grafana.com (thanks [@mjseaman](https://github.com/mjseaman))

## [1.0.7] 2020-12-07

- Added Tags to the list of queryable resources in the AnnotationsQueryEditor (
  thanks [@nazzzzz](https://github.com/nazzzzz))

## [1.0.6] 2020-09-24

- Added a message to the healthcheck success status (thanks [@vladimirdotk](https://github.com/vladimirdotk))
- Added URL option for GitHub Enterprise Users (thanks [@bmike78](https://github.com/bmike78))

## [1.0.5] 2020-09-15

- Added Pull Request ID (Number), URL, and Repository name to pull request responses ( fixes #60 )
- Added the ability to search for all Pull Requests in an organization using the org: search term ( fixes #61 )
- Removed limit from repository list ( fixes #59 )

## [1.0.3] 2020-09-11

- Add the ability to disable time field filtering for pull requests ( fixes #57 )

## [1.0.1] 2020-09-11

- Add the ability to query repositories for variables ( fixes #52 )
- Fix scoped variables for repeating panels ( fixes #51 )
- The default time field for pull requests (Closed At) is now being displayed instead of an empty dropdown

## [1.0.0] 2020-09-10

- Initial release
