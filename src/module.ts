import { DataSourcePlugin } from '@grafana/data';
import { GitHubDataSource } from './DataSource';
import ConfigEditorJsonRender from './views/ConfigEditorJsonRender';
import QueryEditor from './views/QueryEditor';
import type { GitHubQuery } from './types/query';
import type { GitHubDataSourceOptions, GitHubSecureJsonData } from './types/config';

export const plugin = new DataSourcePlugin<
  GitHubDataSource,
  GitHubQuery,
  GitHubDataSourceOptions,
  GitHubSecureJsonData
>(GitHubDataSource)
  .setConfigEditor(ConfigEditorJsonRender)
  .setQueryEditor(QueryEditor);
