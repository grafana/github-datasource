import { TemplateSrv } from '@grafana/runtime';
import { ScopedVars } from '@grafana/data';
import type { GitHubQuery } from './types/query';

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
    updatedQuery = { ...updatedQuery, options: interpolateObject(query.options, t, scoped) };
  }
  return updatedQuery;
};

const interpolateObject = (input: any, t: TemplateSrv, scoped: ScopedVars = {}) => {
  let newOptions = { ...input };
  Object.keys(newOptions).forEach((key) => {
    if (key !== 'refId') {
      const option = newOptions[key];
      if (typeof option === 'string') {
        if (key === 'query') {
          newOptions = { ...newOptions, [key]: replaceVariable(t, option, scoped, 'csv') };
        } else {
          newOptions = { ...newOptions, [key]: replaceVariable(t, option, scoped) };
        }
      } else if (Array.isArray(option)) {
        const replaced = option.map((opt) => interpolateObject(opt, t, scoped));
        newOptions = { ...newOptions, [key]: replaced };
      }
    }
  });
  return newOptions;
};
