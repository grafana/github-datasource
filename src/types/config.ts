import type { DataSourceJsonData } from '@grafana/schema';

export type GitHubLicenseType = 'github-basic' | 'github-enterprise-cloud' | 'github-enterprise-server';

export type GitHubAuthType = 'personal-access-token' | 'github-app';

type GithubCommonOptionsBase<T extends GitHubLicenseType> = {
  githubPlan?: T;
} & DataSourceJsonData

type GitHubDataSourceBasicOptions = GithubCommonOptionsBase<'github-basic'> & {
  githubPlan?: 'github-basic';
  githubUrl: never;
};

type GitHubDataSourceEnterpriseCloudOptions = GithubCommonOptionsBase<'github-enterprise-cloud'> & {
  githubPlan: 'github-enterprise-cloud';
  githubUrl: never;
};

type GitHubDataSourceEnterpriseServerOptions = GithubCommonOptionsBase<'github-enterprise-server'> & {
  githubPlan: 'github-enterprise-server';
  githubUrl: string;
};

type GithubDataSourceCommonOptions = (GitHubDataSourceBasicOptions | GitHubDataSourceEnterpriseCloudOptions | GitHubDataSourceEnterpriseServerOptions)

type GithubDataSourceAuthOptionsBase<T extends GitHubAuthType> = {
  selectedAuthType?: T
}

type GitHubDataSourcePATAuthOptions = GithubDataSourceAuthOptionsBase<'personal-access-token'> & {
  appId: never;
  installationId: string;
};

type GitHubDataSourceGHAppOptions = GithubDataSourceAuthOptionsBase<'github-app'> & {
  appId: string;
  installationId: string;
};

type GithubDataSourceAuthOptions = (GitHubDataSourcePATAuthOptions | GitHubDataSourceGHAppOptions)

export type GitHubDataSourceOptions = GithubDataSourceCommonOptions & GithubDataSourceAuthOptions;

type GitHubSecureJsonDataAuthPAT = {
  accessToken: string;
  privateKey?: string;
};

type GitHubSecureJsonDataAuthGHApp = {
  privateKey: string;
  accessToken?: string;
};

export type GitHubSecureJsonData = GitHubSecureJsonDataAuthPAT | GitHubSecureJsonDataAuthGHApp;
