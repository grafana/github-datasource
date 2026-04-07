import { DataSourcePlugin } from '@grafana/data';
import { GitHubDataSource } from './DataSource';
import ConfigEditor from './views/ConfigEditor';
import QueryEditor from './views/QueryEditor';
import type { GitHubQuery } from './types/query';
import type { GitHubDataSourceOptions, GitHubSecureJsonData } from './types/config';

export const plugin = new DataSourcePlugin<
  GitHubDataSource,
  GitHubQuery,
  GitHubDataSourceOptions,
  GitHubSecureJsonData
>(GitHubDataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
