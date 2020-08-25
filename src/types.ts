import { DataSourceJsonData } from '@grafana/data';

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
