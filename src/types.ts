import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/data';

export interface Label {
  color: string;
  description: string;
  name: string;
}

export interface RepositoryOptions {
  repository?: string;
  owner?: string;
}

export interface GithubDataSourceOptions extends DataSourceJsonData, RepositoryOptions {
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
}

export enum PullRequestTimeField {
  ClosedAt,
  CreatedAt,
  MergedAt,
}

export enum IssueTimeField {
  CreatedAt,
  ClosedAt,
}

export interface IssueFilters {
  assignee?: string;
  createdBy?: string;
  labels?: string[];
  mentioned?: string;
  milestone?: string;
}

export interface ReleasesOptions {}
export interface TagsOptions {}
export interface PullRequestsOptions {
  timeField?: PullRequestTimeField;
  query?: string;
}

export interface CommitsOptions {
  gitRef?: string;
}

export interface ContributorsOptions {
  query?: string;
}

export interface LabelsOptions {
  query?: string;
}

export interface IssuesOptions {
  timeField?: IssueTimeField;
  filters?: IssueFilters;
  query?: string;
}

export interface GitHubQuery extends DataQuery, RepositoryOptions {
  pullRequestsOptions?: PullRequestsOptions;
  releasesOptions?: ReleasesOptions;
  labelsOptions?: LabelsOptions;
  tagsOptions?: TagsOptions;
  commitsOptions?: CommitsOptions;
  issuesOptions?: IssuesOptions;
  contributorsOptions?: ContributorsOptions;
}

export interface GitHubVariableQuery extends GitHubQuery {
  field?: string;
}

export const DefaultQueryType = QueryType.Commits;
