export enum QueryType {
  Commits = 'Commits',
  Issues = 'Issues',
  Contributors = 'Contributors',
  Tags = 'Tags',
  Releases = 'Releases',
  Pull_Requests = 'Pull_Requests',
  Labels = 'Labels',
  Repositories = 'Repositories',
  Organizations = 'Organizations',
  GraphQL = 'GraphQL',
  Milestones = 'Milestones',
  Packages = 'Packages',
  Vulnerabilities = 'Vulnerabilities',
  Projects = 'Projects',
  ProjectItems = 'ProjectItems',
  Stargazers = 'Stargazers',
  Workflows = 'Workflows',
  Workflow_Usage = 'Workflow_Usage',
  Workflow_Runs = 'Workflow_Runs',
}

export const DefaultQueryType = QueryType.Issues;

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
  None,
}

export enum IssueTimeField {
  CreatedAt,
  ClosedAt,
  UpdatedAt,
}

export enum WorkflowsTimeField {
  CreatedAt,
  UpdatedAt,
}

export enum ProjectQueryType {
  ORG = 0,
  USER = 1,
}
