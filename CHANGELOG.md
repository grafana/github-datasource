# Change Log

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

- **Feature** - Github projects - query items, user projects
- **Chore** - Updated build to use go 1.19.3

## [1.2.0] 2022-10-20

- **Feature** - Github projects

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
