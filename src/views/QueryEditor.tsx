import React, { ReactNode, useCallback } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select } from '@grafana/ui';

import { DataSource } from '../DataSource';
import { GithubDataSourceOptions, GitHubQuery, QueryType, DefaultQueryType } from '../types';
import { QueryInlineField } from '../components/Forms';
import { isValid } from '../validation';

import QueryEditorRepository from './QueryEditorRepository';
import QueryEditorReleases from './QueryEditorReleases';
import QueryEditorCommits from './QueryEditorCommits';
import QueryEditorIssues from './QueryEditorIssues';
import QueryEditorPullRequests from './QueryEditorPullRequests';
import QueryEditorTags from './QueryEditorTags';
import QueryEditorContributors from './QueryEditorContributors';
import QueryEditorLabels from './QueryEditorLabels';

interface Props extends QueryEditorProps<DataSource, GitHubQuery, GithubDataSourceOptions> {
  queryTypes?: string[];
} 
export const LeftColumnWidth = 10;
export const RightColumnWidth = 36;

const queryEditors: {
  [key:string]: { component: (props: Props, onChange: (val: any) => void) => ReactNode; optionsKey: string };
} = {
  [QueryType.Labels]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorLabels {...(props.query.labelsOptions || {})} onChange={onChange} />
    ),
    optionsKey: 'labelsOptions',
  },
  [QueryType.Contributors]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorContributors {...(props.query.contributorsOptions || {})} onChange={onChange} />
    ),
    optionsKey: 'contributorsOptions',
  },
  [QueryType.Tags]: {
    component: (props: Props, _: (val: any) => void) => (
      <QueryEditorTags {...(props.query.tagsOptions || {})}  />
    ),
    optionsKey: 'tagsOptions',
  },
  [QueryType.Releases]: {
    component: (props: Props, _: (val: any) => void) => (
      <QueryEditorReleases {...(props.query.releasesOptions || {})} />
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
      <QueryEditorIssues {...(props.query.issuesOptions || {})} onChange={onChange} datasource={props.datasource} owner={props.query.owner || ''} repository={props.query.repository || ''}/>
    ),
    optionsKey: 'issuesOptions',
  },
  [QueryType.Pull_Requests]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorPullRequests {...(props.query.pullRequestsOptions || {})} onChange={onChange} />
    ),
    optionsKey: 'pullRequestsOptions',
  },
};

const queryTypeOptions: Array<SelectableValue<string>> = Object.keys(QueryType)
  .map((v) => {
    return {
      label: v.replace("_", " "),
      value: v,
    };
  });

export default (props: Props) => {
  const onChange = useCallback(
    (value: GitHubQuery) => {
      props.onChange(value);

      if(isValid(value)) {
        props.onRunQuery();
      }
    },
    [props.onChange]);

  const onKeyChange = useCallback(
    (key: string, value: any) => {
        onChange({
        ...props.query,
        [key]: value,
      });
    },
    [onChange]
  );

  const queryEditor = queryEditors[props.query.queryType || DefaultQueryType];
  const queryTypes = props.queryTypes || Object.keys(queryEditors);
  return (
    <>
      <QueryInlineField label="Query Type" tooltip="What resource are you querying for?" labelWidth={LeftColumnWidth}>
        <Select
          width={RightColumnWidth}
          options={queryTypeOptions.filter(v => queryTypes.includes(v.value!))}
          value={props.query.queryType}
          onChange={val => onKeyChange('queryType', val.value || DefaultQueryType)}
        />
      </QueryInlineField>

      <QueryEditorRepository repository={props.query.repository} owner={props.query.owner} onChange={(repo => {
        onChange({
          ...props.query,
          ...repo,
        })
      })} />

      {queryEditor ? (
        queryEditor.component(props, (value: any) => onKeyChange(queryEditor.optionsKey, !!value ? value : undefined))
      ) : (
        <span>Unsupported Query Type</span>
      )}
    </>
  );
};
