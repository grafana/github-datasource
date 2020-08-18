import React, { useCallback } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { QueryType, GitHubQuery, DataSourceOptions } from '../types';
import { QueryInlineField } from '../components/Forms';

type Props = QueryEditorProps<DataSource, GitHubQuery, DataSourceOptions>;

const leftColumnWidth = 12;
const rightColumnWidth = 36;

const queryTypeOptions: SelectableValue<QueryType>[] = Object.keys(QueryType).filter((_, i) => QueryType[i] !== undefined).map((_, i) => {
  return {
    label: `${QueryType[i]}`,
    value: i as QueryType,
  }
});

export default (props: Props) => {
  const onChange = useCallback((key: string, value: any) => {
    console.log(key, value);
    props.onChange({
      ...props.query,
      [key]: value,
      repository: 'grafana',
      owner: 'grafana',
      ref: 'master',
    });

    console.log(props.query);
    props.onRunQuery();
  }, [props.onChange]);

  return (
    <>
      <QueryInlineField label="Resource" tooltip="What resource are you querying for?" labelWidth={leftColumnWidth}>
        <Select width={rightColumnWidth} options={queryTypeOptions} value={props.query.queryType || QueryType.Commits} onChange={(val) => onChange('type', val.value)} />
      </QueryInlineField>
    </>
  );
};
