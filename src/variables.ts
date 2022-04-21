import { GitHubQuery } from './types';
import { TemplateSrv } from '@grafana/runtime';
import { ScopedVars } from '@grafana/data';

export const replaceVariable = (
  t: TemplateSrv,
  value?: string,
  scoped?: ScopedVars,
  format?: string
): string | undefined => {
  return !!value ? t.replace(value, scoped, format) : value;
};

export const replaceVariables = (t: TemplateSrv, query: GitHubQuery, scoped: ScopedVars): GitHubQuery => {
  Object.keys(query).forEach((key) => {
    if (typeof query[key] === 'string') {
      if (key === 'query') {
        query = { ...query, [key]: replaceVariable(t, query[key], scoped, 'csv') };
      } else {
        query = { ...query, [key]: replaceVariable(t, query[key], scoped) };
      }
    }
  });

  if (query.options) {
    let { options } = query;
    Object.keys(options).forEach((key) => {
      if (typeof options[key] === 'string') {
        if (key === 'query') {
          options = { ...options, [key]: replaceVariable(t, options[key], scoped, 'csv') };
        } else {
          options = { ...options, [key]: replaceVariable(t, options[key], scoped) };
        }
      }
    });
    query.options = options;
  }

  return query;
};
