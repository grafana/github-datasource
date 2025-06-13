import { map, Observable } from 'rxjs';
import { TemplateSrv } from '@grafana/runtime';
import { ScopedVars, CustomVariableSupport, DataQueryResponse, DataQueryRequest, DataFrameView } from '@grafana/data';
import { GitHubDataSource } from './DataSource';
import VariableQueryEditor from 'views/VariableQueryEditor';
import type { GitHubQuery, GitHubVariableQuery } from './types/query';

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

export class GithubVariableSupport extends CustomVariableSupport<GitHubDataSource, GitHubVariableQuery, GitHubQuery> {
  constructor(private readonly datasource: GitHubDataSource) {
    super();
    this.datasource = datasource;
    this.query = this.query.bind(this);
  }
  editor = VariableQueryEditor;
  query(request: DataQueryRequest<GitHubVariableQuery>): Observable<DataQueryResponse> {
    let query = { ...request?.targets[0], refId: 'metricFindQuery' };
    return this.datasource
      .query({ ...request, targets: [query] })
      .pipe(map((response) => ({ ...response, data: response.data || [] })))
      .pipe(map((response) => queryResponseToVariablesFrame(query, response)));
  }
}

const queryResponseToVariablesFrame = (query: GitHubVariableQuery, response: DataQueryResponse) => {
  if (response?.data?.length < 1) {
    return { ...response, data: [] };
  }
  const view = new DataFrameView(response.data[0] || {});
  const data = view.map((item) => {
    const value = item[query.key || ''] || item[query.field || 'name'];
    const text = item[query.field || 'name'] || value;
    return { value, text };
  });
  return { ...response, data };
};
