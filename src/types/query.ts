import { type DataQuery } from '@grafana/schema';
import { PullRequestTimeField, IssueTimeField, WorkflowsTimeField, PackageType, ProjectQueryType, QueryTypes } from '../constants';
import type { Filter } from 'components/Filters';

export type QueryType = typeof QueryTypes[number]

export type Options = { [index: string]: any; }

export type RepositoryOptions = { repository?: string; owner?: string; }

type BaseQuery<T extends QueryType, O extends Options> = { queryType: T; options?: O; } & RepositoryOptions & DataQuery;

//#region Code_Scanning Query
export type CodeScanningOptions = Options & {
  gitRef?: string;
  state?: string
}
type Code_ScanningQuery = BaseQuery<'Code_Scanning', CodeScanningOptions>
//#endregion

//#region Commits Query
export type CommitsOptions = Options & {
  gitRef?: string;
}
type CommitsQuery = BaseQuery<'Commits', CommitsOptions>
//#endregion

//#region Commit_Files Query
export type CommitFilesOptions = Options & {
  commitSha?: string;
}
type Commit_FilesQuery = BaseQuery<'Commit_Files', CommitFilesOptions>
//#endregion

//#region Issues Query
export type IssuesOptions = Options & {
  timeField?: IssueTimeField;
  query?: string;
}
type IssuesQuery = BaseQuery<'Issues', IssuesOptions>
//#endregion

//#region Contributors Query
export type ContributorsOptions = Options & {
  query?: string;
}
type ContributorsQuery = BaseQuery<'Contributors', ContributorsOptions>
//#endregion

//#region Tags Query
export type TagsOptions = Options & {}
type TagsQuery = BaseQuery<'Tags', TagsOptions>
//#endregion

//#region Releases Query
export type ReleasesOptions = Options & {}
type ReleasesQuery = BaseQuery<'Releases', ReleasesOptions>
//#endregion

//#region Pull_Requests Query
export type Pull_RequestsOptions = Options & {
  timeField?: PullRequestTimeField;
  query?: string;
}
type Pull_RequestsQuery = BaseQuery<'Pull_Requests', Pull_RequestsOptions>
//#endregion

//#region Pull_Request_Files Query
export type Pull_Request_FilesOption = Options & {
  prNumber?: number;
}
type Pull_Request_FilesQuery = BaseQuery<'Pull_Request_Files', Pull_Request_FilesOption>
//#endregion

//#region Pull_Request_Reviews Query
export type PullRequestReviewsOptions = Options & {
  timeField?: PullRequestTimeField;
  query?: string;
}
type Pull_Request_ReviewsQuery = BaseQuery<'Pull_Request_Reviews', PullRequestReviewsOptions>
//#endregion

//#region Labels Query
export type LabelsOptions = Options & {
  query?: string;
}
type LabelsQuery = BaseQuery<'Labels', LabelsOptions>
//#endregion

//#region Repositories Query
type RepositoriesQuery = BaseQuery<'Repositories', {}>
//#endregion

//#region Organizations Query
type OrganizationsQuery = BaseQuery<'Organizations', {}>
//#endregion

//#region GraphQL Query
type GraphQLQuery = BaseQuery<'GraphQL', {}>
//#endregion

//#region Milestones Query
export type MilestonesOptions = Options & {
  query?: string;
}
type MilestonesQuery = BaseQuery<'Milestones', MilestonesOptions>
//#endregion

//#region Packages Query
export type PackagesOptions = Options & {
  names?: string;
  packageType?: PackageType;
}
type PackagesQuery = BaseQuery<'Packages', PackagesOptions>
//#endregion

//#region Vulnerabilities Query
type VulnerabilitiesQuery = BaseQuery<'Vulnerabilities', {}>
//#endregion

//#region Projects Query
export type ProjectsOptions = Options & {
  organization?: string;
  number?: number | string;
  user?: string;
  kind?: ProjectQueryType;
  filters?: Filter[];
}
type ProjectsQuery = BaseQuery<'Projects', ProjectsOptions>
//#endregion

//#region ProjectItems Query
type ProjectItemsQuery = BaseQuery<'ProjectItems', {}>
//#endregion

//#region Stargazers Query
type StargazersQuery = BaseQuery<'Stargazers', {}>
//#endregion

//#region Workflow_Runs Query
export type WorkflowRunsOptions = Options & {
  workflowID?: string;
  branch?: string;
}
type Workflow_RunsQuery = BaseQuery<'Workflow_Runs', WorkflowRunsOptions>
//#endregion

//#region Workflow_Usage Query
export type WorkflowUsageOptions = Options & {
  workflowID?: number;
}
type Workflow_UsageQuery = BaseQuery<'Workflow_Usage', WorkflowUsageOptions>
//#endregion

//#region Workflows Query
export type WorkflowsOptions = Options & {
  timeField?: WorkflowsTimeField;
  query?: string;
}
type WorkflowsQuery = BaseQuery<'Workflows', WorkflowsOptions>
//#endregion

//#region Deployments Query
export type DeploymentsOptions = Options & {
  sha?: string;
  gitRef?: string;
  task?: string;
  environment?: string;
}
type DeploymentsQuery = BaseQuery<'Deployments', DeploymentsOptions>
//#endregion

export type GitHubQuery =
  Code_ScanningQuery |
  CommitsQuery |
  Commit_FilesQuery |
  IssuesQuery |
  ContributorsQuery |
  TagsQuery |
  ReleasesQuery |
  Pull_RequestsQuery |
  Pull_Request_ReviewsQuery |
  Pull_Request_FilesQuery |
  LabelsQuery |
  RepositoriesQuery |
  OrganizationsQuery |
  GraphQLQuery |
  MilestonesQuery |
  PackagesQuery |
  VulnerabilitiesQuery |
  ProjectsQuery |
  ProjectItemsQuery |
  StargazersQuery |
  Workflow_RunsQuery |
  Workflow_UsageQuery |
  WorkflowsQuery |
  DeploymentsQuery

export type GitHubVariableQuery = { key?: string; field?: string; } & GitHubQuery

export type GitHubAnnotationQuery = { timeField?: string; } & GitHubVariableQuery
