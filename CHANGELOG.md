# Change Log

## Entries

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

- Fixed issue where some time values were being renderred incorrectly

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
