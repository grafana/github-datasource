import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import ConfigEditor from './views/ConfigEditorContainer';
import QueryEditor from './views/QueryEditor';
import { GitHubQuery, DataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, GitHubQuery, DataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
