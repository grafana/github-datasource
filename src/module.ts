import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import ConfigEditor from './components/ConfigEditor';
import QueryEditor from './components/QueryEditor';
import { GitHubQuery, DataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, GitHubQuery, DataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
