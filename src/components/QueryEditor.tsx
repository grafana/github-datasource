import React from 'react';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../DataSource';
import { GitHubQuery, DataSourceOptions } from '../types';

type Props = QueryEditorProps<DataSource, GitHubQuery, DataSourceOptions>;

export default (_: Props) => {
  return <></>
}
