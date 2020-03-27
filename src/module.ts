import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor, QueryEditor, MetaInspector } from './components';
import { SheetsQuery, SheetsSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, SheetsQuery, SheetsSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setMetadataInspector(MetaInspector);
