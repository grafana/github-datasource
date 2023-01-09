import { DataSourcePlugin } from '@grafana/data';
import { GithubDataSource } from './DataSource';
import ConfigEditor from './views/ConfigEditor';
import QueryEditor from './views/QueryEditor';
import VariableQueryEditor from './views/VariableQueryEditor';
import { GitHubQuery, GithubDataSourceOptions, GithubSecureJsonData } from './types';

export const plugin = new DataSourcePlugin<
  GithubDataSource,
  GitHubQuery,
  GithubDataSourceOptions,
  GithubSecureJsonData
>(GithubDataSource)
  .setConfigEditor(ConfigEditor)
  .setVariableQueryEditor(VariableQueryEditor)
  .setQueryEditor(QueryEditor);
