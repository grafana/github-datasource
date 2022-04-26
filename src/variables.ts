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
  let updatedQuery = { ...query };
  updatedQuery = interpolateObject(updatedQuery, t, scoped);
  if (query.options) {
    updatedQuery = { ...query, options: interpolateObject(query.options, t, scoped) };
  }
  return updatedQuery;
};

const interpolateObject = (input: any, t: TemplateSrv, scoped: ScopedVars = {}) => {
  let newOptions = { ...input };
  Object.keys(newOptions).forEach((key) => {
    if (key !== 'refId') {
      if (typeof newOptions[key] === 'string') {
        if (key === 'query') {
          newOptions = { ...newOptions, [key]: replaceVariable(t, newOptions[key], scoped, 'csv') };
        } else {
          newOptions = { ...newOptions, [key]: replaceVariable(t, newOptions[key], scoped) };
        }
      }
    }
  });
  return newOptions;
};
