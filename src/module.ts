import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import ConfigEditor from './views/ConfigEditorContainer';
import QueryEditor from './views/QueryEditor';
import { DataSourceOptions } from './types';
import { GitHubQuery } from './query';

export const plugin = new DataSourcePlugin<DataSource, GitHubQuery, DataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
