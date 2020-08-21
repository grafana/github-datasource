import { DataQuery } from '@grafana/data';

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

export interface IssueFilters {
  assignee?: string;
  createdBy?: string;
  labels?: string[];
  mentioned?: string;
  milestone?: string;
}

export interface RepositoryOptions {
  repository?: string;
  owner?: string;
}

export interface ReleasesOptions extends RepositoryOptions {}
export interface TagsOptions extends RepositoryOptions {}
export interface PullRequestsOptions extends RepositoryOptions {
  timeField?: PullRequestTimeField;
}
export interface ContributorsOptions extends RepositoryOptions {
  gitRef?: string;
}
export interface CommitsOptions extends RepositoryOptions {
  gitRef?: string;
}
export interface IssuesOptions extends RepositoryOptions {
  filters?: IssueFilters;
}

export interface GitHubQuery extends DataQuery {
  type: QueryType;
  pullRequestsOptions?: PullRequestsOptions;
  releasesOptions?: ReleasesOptions;
  tagsOptions?: TagsOptions;
  commitsOptions?: CommitsOptions;
  issuesOptions?: IssuesOptions;
}
