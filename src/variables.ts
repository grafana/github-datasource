import { GitHubQuery } from './types';
import { TemplateSrv } from '@grafana/runtime';

export const ReplaceVariable = (t: TemplateSrv, value?: string): string | undefined => {
  return !!value ? t.replace(value) : value;
};

export const ReplaceVariables = (t: TemplateSrv, query: GitHubQuery): GitHubQuery => {
  return {
    ...query,
    owner: ReplaceVariable(t, query.owner),
    repository: ReplaceVariable(t, query.repository),
  };
};
