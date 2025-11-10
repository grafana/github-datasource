import type { DataSourceJsonData } from '@grafana/data';

export type GitHubLicenseType = 'github-basic' | 'github-enterprise-cloud' | 'github-enterprise-server';

export type GitHubAuthType = 'personal-access-token' | 'github-app';

export type GitHubDataSourceOptions = {
  githubPlan?: GitHubLicenseType;
  githubUrl?: string;
  selectedAuthType?: GitHubAuthType;
  appId?: string;
  installationId?: string;
} & DataSourceJsonData;

export type GitHubSecureJsonDataKeys =
  | 'accessToken' // accessToken is set if the user is using a Personal Access Token to connect to GitHub
  | 'privateKey'; // privateKey is set if the user is using a GitHub App to connect to GitHub

export type GitHubSecureJsonData = Partial<Record<GitHubSecureJsonDataKeys, string>>;
