import { z } from 'zod';
import type { DataSourceJsonData } from '@grafana/schema';

//#region jsonData

//#region --- License / Plan schemas ---

const GitHubLicenseTypeSchema = z.enum(['github-basic', 'github-enterprise-cloud', 'github-enterprise-server']).describe('The GitHub license/plan type');
export type GitHubLicenseType = z.infer<typeof GitHubLicenseTypeSchema>;

const GitHubAuthTypeSchema = z.enum(['personal-access-token', 'github-app']).describe('The GitHub authentication method');
export type GitHubAuthType = z.infer<typeof GitHubAuthTypeSchema>;

//#endregion

//#region --- Plan option schemas ---

const GitHubDataSourceBasicOptionsSchema = z.object({
  githubPlan: z.literal('github-basic').optional().describe('GitHub plan type (basic)'),
  githubUrl: z.never().describe('Not applicable for GitHub basic plan'),
}).describe('Configuration for GitHub basic plan');

const GitHubDataSourceEnterpriseCloudOptionsSchema = z.object({
  githubPlan: z.literal('github-enterprise-cloud').describe('GitHub plan type (Enterprise Cloud)'),
  githubUrl: z.never().describe('Not applicable for GitHub Enterprise Cloud'),
}).describe('Configuration for GitHub Enterprise Cloud plan');

const GitHubDataSourceEnterpriseServerOptionsSchema = z.object({
  githubPlan: z.literal('github-enterprise-server').describe('GitHub plan type (Enterprise Server)'),
  githubUrl: z.string().describe('The URL of the GitHub Enterprise Server instance'),
}).describe('Configuration for GitHub Enterprise Server plan');

const GithubDataSourceCommonOptionsSchema = GitHubDataSourceBasicOptionsSchema
  .or(GitHubDataSourceEnterpriseCloudOptionsSchema)
  .or(GitHubDataSourceEnterpriseServerOptionsSchema);

//#endregion

//#region --- Auth option schemas ---

const GitHubDataSourcePATAuthOptionsSchema = z.object({
  selectedAuthType: z.literal('personal-access-token').optional().describe('Authentication type (Personal Access Token)'),
  appId: z.never().describe('Not applicable for PAT authentication'),
  installationId: z.never().describe('Not applicable for PAT authentication'),
}).describe('Authentication options for Personal Access Token');

const GitHubDataSourceGHAppOptionsSchema = z.object({
  selectedAuthType: z.literal('github-app').describe('Authentication type (GitHub App)'),
  appId: z.string().describe('The GitHub App ID'),
  installationId: z.string().describe('The GitHub App installation ID'),
}).describe('Authentication options for GitHub App');

const GithubDataSourceAuthOptionsSchema = GitHubDataSourcePATAuthOptionsSchema
  .or(GitHubDataSourceGHAppOptionsSchema)

//#endregion

export const GitHubDataSourceOptionsSchema = z.intersection(GithubDataSourceCommonOptionsSchema, GithubDataSourceAuthOptionsSchema).describe('GitHub data source configuration options (jsonData)');

export type GitHubDataSourceOptions = z.infer<typeof GitHubDataSourceOptionsSchema> & DataSourceJsonData;

//#endregion

//#region secureJsonData

//#region --- Secure JSON data schemas ---

const GitHubSecureJsonDataAuthPATSchema = z.object({
  accessToken: z.string().describe('Personal access token for GitHub API authentication'),
  privateKey: z.never().describe('Not applicable for PAT authentication'),
}).describe('Secure data for Personal Access Token authentication');

const GitHubSecureJsonDataAuthGHAppSchema = z.object({
  accessToken: z.never().describe('Not applicable for GitHub App authentication'),
  privateKey: z.string().describe('Private key for GitHub App authentication (PEM format)'),
}).describe('Secure data for GitHub App authentication');

//#endregion

export const GitHubSecureJsonDataSchema = GitHubSecureJsonDataAuthPATSchema
  .or(GitHubSecureJsonDataAuthGHAppSchema)
  .describe('Secure JSON data for GitHub data source authentication (secureJsonData)');

export type GitHubSecureJsonData = z.infer<typeof GitHubSecureJsonDataSchema>;

//#endregion
