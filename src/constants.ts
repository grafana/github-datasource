import { QueryType } from './types/query'

export const QueryTypes = [
  'Code_Scanning',
  'Commits',
  'Issues',
  'Contributors',
  'Tags',
  'Releases',
  'Pull_Requests',
  'Pull_Request_Reviews',
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
  'Deployments'
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
