import { GitHubQuery } from './types';
import { TemplateSrv } from '@grafana/runtime';

export const replaceVariable = (t: TemplateSrv, value?: string): string | undefined => {
  return !!value ? t.replace(value) : value;
};

export const replaceVariables = (t: TemplateSrv, query: GitHubQuery): GitHubQuery => {
  Object.keys(query).forEach(key => {
    if (typeof query[key] === 'string') {
      query[key] = replaceVariable(t, query[key]);
    }
  });

  if (query.options) {
    const { options } = query;
    Object.keys(options).forEach(key => {
      if (typeof options[key] === 'string') {
        options[key] = replaceVariable(t, options[key]);
      }
    });
    query.options = options;
  }

  return query;
};
