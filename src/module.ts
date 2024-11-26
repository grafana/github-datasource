import { DataSourcePlugin } from '@grafana/data';
import { GitHubDataSource } from './DataSource';
import ConfigEditor from './views/ConfigEditor';
import QueryEditor from './views/QueryEditor';
import VariableQueryEditor from './views/VariableQueryEditor';
import type { GitHubQuery } from './types';
import type { GitHubDataSourceOptions, GitHubSecureJsonData } from './types/config';

export const plugin = new DataSourcePlugin<
  GitHubDataSource,
  GitHubQuery,
  GitHubDataSourceOptions,
  GitHubSecureJsonData
>(GitHubDataSource)
  .setConfigEditor(ConfigEditor)
  .setVariableQueryEditor(VariableQueryEditor)
  .setQueryEditor(QueryEditor);
