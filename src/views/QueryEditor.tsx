import React, { ReactNode, useCallback } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { DataSourceOptions } from '../types';
import { GitHubQuery, QueryType } from '../query';
import { QueryInlineField } from '../components/Forms';

import QueryEditorReleases from './QueryEditorReleases';
import QueryEditorCommits from './QueryEditorCommits';
import QueryEditorIssues from './QueryEditorIssues';
import QueryEditorPullRequests from './QueryEditorPullRequests';

export type Props = QueryEditorProps<DataSource, GitHubQuery, DataSourceOptions>;

const queryEditors: {
  [key: number]: { component: (props: Props, onChange: (val: any) => void) => ReactNode; optionsKey: string };
} = {
  [QueryType.Releases]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorReleases {...(props.query.releasesOptions || {})} onChange={onChange} />
    ),
    optionsKey: 'releasesOptions',
  },
  [QueryType.Commits]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorCommits {...(props.query.commitsOptions || {})} onChange={onChange} />
    ),
    optionsKey: 'commitsOptions',
  },
  [QueryType.Issues]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorIssues {...(props.query.issuesOptions || {})} onChange={onChange} datasource={props.datasource} />
    ),
    optionsKey: 'issuesOptions',
  },
  [QueryType.PullRequests]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorPullRequests {...(props.query.pullRequestsOptions || {})} onChange={onChange} />
    ),
    optionsKey: 'pullRequestsOptions',
  },
};

export const LeftColumnWidth = 12;
export const RightColumnWidth = 36;

const queryTypeOptions: Array<SelectableValue<QueryType>> = Object.keys(QueryType)
  .filter((_, i) => QueryType[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${QueryType[i]}`,
      value: i as QueryType,
    };
  });

export default (props: Props) => {
  const onChange = useCallback(
    (key: string, value: any) => {
      props.onChange({
        ...props.query,
        [key]: value,
      });

      props.onRunQuery();
    },
    [props.onChange]
  );

  const queryEditor = queryEditors[props.query.type];

  return (
    <>
      <QueryInlineField label="Query Type" tooltip="What resource are you querying for?" labelWidth={LeftColumnWidth}>
        <Select
          width={RightColumnWidth}
          options={queryTypeOptions}
          value={props.query.type}
          onChange={val => onChange('type', val.value)}
        />
      </QueryInlineField>
      {queryEditor ? (
        queryEditor.component(props, (value: any) => onChange(queryEditor.optionsKey, value))
      ) : (
        <span>Unsupported Query Type</span>
      )}
    </>
  );
};
