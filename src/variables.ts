import { GitHubQuery } from './types';
import { TemplateSrv } from '@grafana/runtime';
import { ScopedVars } from '@grafana/data';

export const replaceVariable = (t: TemplateSrv, value?: string, scoped: ScopedVars): string | undefined => {
  return !!value ? t.replace(value, scoped) : value;
};

export const replaceVariables = (t: TemplateSrv, query: GitHubQuery, scoped: ScopedVars): GitHubQuery => {
  Object.keys(query).forEach((key) => {
    if (typeof query[key] === 'string') {
      query = { ...query, [key]: replaceVariable(t, query[key], scoped) };
    }
  });

  if (query.options) {
    let { options } = query;
    Object.keys(options).forEach((key) => {
      if (typeof options[key] === 'string') {
        options = {...options, [key]: replaceVariable(t, options[key], scoped)}
      }
    });
    query.options = options;
  }

  return query;
};
