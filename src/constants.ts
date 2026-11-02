import type { PullRequestField, QueryType } from './types/query'

export const QueryTypes = [
  'Check_Runs',
  'Code_Scanning',
  'Commits',
  'Commit_Files',
  'Issues',
  'Contributors',
  'Tags',
  'Branches',
  'Releases',
  'Pull_Requests',
  'Pull_Request_Reviews',
  'Pull_Request_Files',
  'Labels',
  'Repositories',
  'Organizations',
  'GraphQL',
  'Milestones',
  'Packages',
  'Vulnerabilities',
  'Projects',
  'ProjectItems',
  'Stargazers',
  'Workflows',
  'Workflow_Usage',
  'Workflow_Runs',
  'Deployments',
] as const;


export const DefaultQueryType: QueryType = "Issues";

export enum PackageType {
  NPM = 'NPM',
  RUBYGEMS = 'RUBYGEMS',
  MAVEN = 'MAVEN',
  DOCKER = 'DOCKER',
  DEBIAN = 'DEBIAN',
  NUGET = 'NUGET',
  PYPI = 'PYPI',
}

export enum PullRequestTimeField {
  ClosedAt,
  CreatedAt,
  MergedAt,
  UpdatedAt,
  None,
}

// The additional fields of a Pull_Requests query, in the order they are shown and reported
export const PullRequestFields: PullRequestField[] = ['diffstat', 'mergeable', 'author', 'merged_by', 'labels'];

// Every additional field is fetched when fields is unset, or when it lists them all
export const hasAllPullRequestFields = (fields?: PullRequestField[]) =>
  fields === undefined || PullRequestFields.every((f) => fields.includes(f));

export enum IssueTimeField {
  CreatedAt,
  ClosedAt,
  UpdatedAt,
}

export enum WorkflowsTimeField {
  None,
  CreatedAt,
  UpdatedAt,
}

export enum ProjectQueryType {
  ORG = 0,
  USER = 1,
}
