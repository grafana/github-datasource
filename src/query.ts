import { DataQuery } from '@grafana/data';
import { RepositoryOptions } from './types';

export enum QueryType {
  Commits,
  Issues,
  Contributors,
  Tags,
  Releases,
  PullRequests,
  GraphQL,
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

export interface ReleasesOptions extends RepositoryOptions {}
export interface TagsOptions extends RepositoryOptions {}
export interface PullRequestsOptions extends RepositoryOptions {
  timeField?: PullRequestTimeField;
  query?: string;
}
export interface ContributorsOptions extends RepositoryOptions {
  gitRef?: string;
}
export interface CommitsOptions extends RepositoryOptions {
  gitRef?: string;
}
export interface IssuesOptions extends RepositoryOptions {
  timeField?: IssueTimeField;
  filters?: IssueFilters;
  query?: string;
}

export interface GitHubQuery extends DataQuery {
  type: QueryType;
  pullRequestsOptions?: PullRequestsOptions;
  releasesOptions?: ReleasesOptions;
  tagsOptions?: TagsOptions;
  commitsOptions?: CommitsOptions;
  issuesOptions?: IssuesOptions;
}
