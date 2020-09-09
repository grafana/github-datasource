import { GitHubQuery } from './types';
import { TemplateSrv } from '@grafana/runtime';

export const ReplaceVariable = (t: TemplateSrv, value?: string): string | undefined => {
  return !!value ? t.replace(value) : value;
};

export const ReplaceVariables = (t: TemplateSrv, query: GitHubQuery): GitHubQuery => {
  Object.keys(query).map(key => {
    if(typeof query[key] === 'string') {
      query[key] = ReplaceVariable(t, query[key])
    }
  });

  if(query.options) {
    const { options } = query;
    Object.keys(options).map(key => {
      if(typeof options[key] === 'string') {
        options[key] = ReplaceVariable(t, options[key]);
      }
    });
    query.options = options;
  }

  return query;
};
