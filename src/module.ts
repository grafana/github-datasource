import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import ConfigEditor from './views/ConfigEditor';
import QueryEditor from './views/QueryEditor';
import { GithubDataSourceOptions, GithubSecureJsonData } from './types';
import { GitHubQuery } from './query';

export const plugin = new DataSourcePlugin<DataSource, GitHubQuery, GithubDataSourceOptions, GithubSecureJsonData>(
  DataSource
)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
