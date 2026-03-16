import { z } from 'zod';
import type { DataSourceJsonData } from '@grafana/schema';

//#region jsonData

//#region --- License / Plan schemas ---

const GitHubLicenseTypeSchema = z.enum(['github-basic', 'github-enterprise-cloud', 'github-enterprise-server']);
export type GitHubLicenseType = z.infer<typeof GitHubLicenseTypeSchema>;

const GitHubAuthTypeSchema = z.enum(['personal-access-token', 'github-app']);
export type GitHubAuthType = z.infer<typeof GitHubAuthTypeSchema>;

//#endregion

//#region --- Plan option schemas ---

const GitHubDataSourceBasicOptionsSchema = z.object({
  githubPlan: z.literal('github-basic').optional(),
  githubUrl: z.never(),
});

const GitHubDataSourceEnterpriseCloudOptionsSchema = z.object({
  githubPlan: z.literal('github-enterprise-cloud'),
  githubUrl: z.never(),
});

const GitHubDataSourceEnterpriseServerOptionsSchema = z.object({
  githubPlan: z.literal('github-enterprise-server'),
  githubUrl: z.string(),
});

const GithubDataSourceCommonOptionsSchema = GitHubDataSourceBasicOptionsSchema
  .or(GitHubDataSourceEnterpriseCloudOptionsSchema)
  .or(GitHubDataSourceEnterpriseServerOptionsSchema);

//#endregion

//#region --- Auth option schemas ---

const GitHubDataSourcePATAuthOptionsSchema = z.object({
  selectedAuthType: z.literal('personal-access-token').optional(),
  appId: z.never(),
  installationId: z.never(),
});

const GitHubDataSourceGHAppOptionsSchema = z.object({
  selectedAuthType: z.literal('github-app'),
  appId: z.string(),
  installationId: z.string(),
});

const GithubDataSourceAuthOptionsSchema = GitHubDataSourcePATAuthOptionsSchema
  .or(GitHubDataSourceGHAppOptionsSchema)

//#endregion

export const GitHubDataSourceOptionsSchema = z.intersection(GithubDataSourceCommonOptionsSchema, GithubDataSourceAuthOptionsSchema);

export type GitHubDataSourceOptions = z.infer<typeof GitHubDataSourceOptionsSchema> & DataSourceJsonData;

//#endregion

//#region secureJsonData

//#region --- Secure JSON data schemas ---

const GitHubSecureJsonDataAuthPATSchema = z.object({
  accessToken: z.string(),
  privateKey: z.never(),
});

const GitHubSecureJsonDataAuthGHAppSchema = z.object({
  accessToken: z.never(),
  privateKey: z.string(),
});

//#endregion

export const GitHubSecureJsonDataSchema = GitHubSecureJsonDataAuthPATSchema
  .or(GitHubSecureJsonDataAuthGHAppSchema)

export type GitHubSecureJsonData = z.infer<typeof GitHubSecureJsonDataSchema>;

//#endregion
