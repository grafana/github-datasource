import { DataQuery, DataSourceJsonData } from '@grafana/data';
import { Filter } from 'components/Filters';

export interface Label {
  color: string;
  description: string;
  name: string;
}

export interface RepositoryOptions {
  repository?: string;
  owner?: string;
}

export interface GithubEnterpriseOptions {
  githubUrl?: string;
}

export interface GithubDataSourceOptions extends DataSourceJsonData, RepositoryOptions, GithubEnterpriseOptions {
  // Any global settings
}

export interface GithubSecureJsonData {
  // accessToken is set if the user is using a Personal Access Token to connect to GitHub
  accessToken?: string;
}

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
}

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
}

export interface Indexable {
  [index: string]: any;
}

export interface ReleasesOptions extends Indexable {}
export interface TagsOptions extends Indexable {}
export interface PullRequestsOptions extends Indexable {
  timeField?: PullRequestTimeField;
  query?: string;
}

export interface CommitsOptions extends Indexable {
  gitRef?: string;
}

export interface ContributorsOptions extends Indexable {
  query?: string;
}

export interface LabelsOptions extends Indexable {
  query?: string;
}

export interface IssuesOptions extends Indexable {
  timeField?: IssueTimeField;
  query?: string;
}

export interface PackagesOptions extends Indexable {
  names?: string;
  packageType?: PackageType;
}

export interface MilestonesOptions extends Indexable {
  query?: string;
}

export interface ProjectsOptions extends Indexable {
  organization?: string;
  number?: number | string;
  user?: string;
  kind?: ProjectQueryType;
  filters?: Filter[];
}

export enum ProjectQueryType {
  ORG = 0,
  USER = 1,
}

export interface GitHubQuery extends Indexable, DataQuery, RepositoryOptions {
  options?:
    | PullRequestsOptions
    | ReleasesOptions
    | LabelsOptions
    | TagsOptions
    | CommitsOptions
    | IssuesOptions
    | ContributorsOptions
    | ProjectsOptions;
}

export interface GitHubVariableQuery extends GitHubQuery {
  key?: string;
  field?: string;
}

export interface GitHubAnnotationQuery extends GitHubVariableQuery {
  timeField?: string;
}

export const DefaultQueryType = QueryType.Issues;
