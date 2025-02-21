import { PullRequestTimeField, IssueTimeField, WorkflowsTimeField, PackageType, ProjectQueryType } from '../constants';
import type { DataQuery } from '@grafana/schema';
import type { Filter } from 'components/Filters';

export interface RepositoryOptions {
  repository?: string;
  owner?: string;
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
    | ProjectsOptions
    | WorkflowsOptions
    | WorkflowUsageOptions
    | WorkflowRunsOptions;
}

export interface Label {
  color: string;
  description: string;
  name: string;
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

export interface WorkflowsOptions extends Indexable {
  timeField?: WorkflowsTimeField;
  query?: string;
}

export interface WorkflowUsageOptions extends Indexable {
  workflowID?: number;
}

export interface WorkflowRunsOptions extends Indexable {
  workflowID?: string;
  branch?: string;
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

export interface GitHubVariableQuery extends GitHubQuery {
  key?: string;
  field?: string;
}

export interface GitHubAnnotationQuery extends GitHubVariableQuery {
  timeField?: string;
}
