import React, { ReactNode, useCallback } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Select, InlineField } from '@grafana/ui';

import { GitHubDataSource } from '../DataSource';
import { isValid } from '../validation';
import { components } from './../components/selectors';

import QueryEditorRepository from './QueryEditorRepository';
import QueryEditorReleases from './QueryEditorReleases';
import QueryEditorCommits from './QueryEditorCommits';
import QueryEditorIssues from './QueryEditorIssues';
import QueryEditorMilestones from './QueryEditorMilestones';
import QueryEditorPullRequests from './QueryEditorPullRequests';
import QueryEditorTags from './QueryEditorTags';
import QueryEditorContributors from './QueryEditorContributors';
import QueryEditorLabels from './QueryEditorLabels';
import QueryEditorPackages from './QueryEditorPackages';
import QueryEditorVulnerabilities from './QueryEditorVulnerabilities';
import QueryEditorProjects from './QueryEditorProjects';
import QueryEditorWorkflows from './QueryEditorWorkflows';
import QueryEditorWorkflowUsage from './QueryEditorWorkflowUsage';
import QueryEditorWorkflowRuns from './QueryEditorWorkflowRuns';
import { QueryType, DefaultQueryType } from '../constants';
import type { GitHubQuery } from '../types/query';
import type { GitHubDataSourceOptions } from '../types/config';

interface Props extends QueryEditorProps<GitHubDataSource, GitHubQuery, GitHubDataSourceOptions> {
  queryTypes?: string[];
}
export const LeftColumnWidth = 10;
export const RightColumnWidth = 36;

/* eslint-disable react/display-name */
const queryEditors: {
  [key: string]: { component: (props: Props, onChange: (val: any) => void) => ReactNode };
} = {
  [QueryType.Repositories]: {
    component: (_: Props, onChange: (val: any) => void) => <></>,
  },
  [QueryType.Labels]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorLabels {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Contributors]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorContributors {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Tags]: {
    component: (props: Props, _: (val: any) => void) => <QueryEditorTags {...(props.query.options || {})} />,
  },
  [QueryType.Releases]: {
    component: (props: Props, _: (val: any) => void) => <QueryEditorReleases {...(props.query.options || {})} />,
  },
  [QueryType.Commits]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorCommits {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Milestones]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorMilestones {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Issues]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorIssues {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Packages]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorPackages {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Pull_Requests]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorPullRequests {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Vulnerabilities]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorVulnerabilities {...(props.query.options || {})} />
    ),
  },
  [QueryType.Projects]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorProjects {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Stargazers]: {
    component: (_: Props, onChange: (val: any) => void) => <></>,
  },
  [QueryType.Workflows]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorWorkflows {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Workflow_Usage]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorWorkflowUsage {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  [QueryType.Workflow_Runs]: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorWorkflowRuns {...(props.query.options || {})} onChange={onChange} />
    ),
  },
};

/* eslint-enable react/display-name */

const queryTypeOptions: Array<SelectableValue<string>> = Object.keys(QueryType).map((v) => {
  return {
    label: v.replace(/_/gi, ' '),
    value: v,
  };
});

const QueryEditor = (props: Props) => {
  const onChange = useCallback(
    (value: GitHubQuery) => {
      props.onChange(value);

      if (isValid(value)) {
        props.onRunQuery();
      }
    },
    [props]
  );

  const onKeyChange = useCallback(
    (key: string, value: any) => {
      onChange({
        ...props.query,
        [key]: value,
      });
    },
    [onChange, props.query]
  );

  const queryEditor = queryEditors[props.query.queryType || DefaultQueryType];
  const queryTypes = props.queryTypes || Object.keys(queryEditors);

  return (
    <>
      <InlineField label="Query Type" tooltip="What resource are you querying for?" labelWidth={LeftColumnWidth * 2}>
        <div aria-label={components.QueryEditor.QueryType.container.ariaLabel}>
          <Select
            menuShouldPortal={true}
            width={RightColumnWidth}
            options={queryTypeOptions.filter((v) => queryTypes.includes(v.value!))}
            value={props.query.queryType}
            onChange={(val) => onKeyChange('queryType', val.value || DefaultQueryType)}
          />
        </div>
      </InlineField>

      {hasRepo(props.query.queryType) && (
        <QueryEditorRepository
          repository={props.query.repository}
          owner={props.query.owner}
          onChange={(repo) => {
            onChange({
              ...props.query,
              ...repo,
            });
          }}
        ></QueryEditorRepository>
      )}

      {queryEditor ? (
        queryEditor.component(props, (value: any) => onKeyChange('options', !!value ? value : undefined))
      ) : (
        <span>Unsupported Query Type</span>
      )}
    </>
  );
};

const nonRepoTypes = [QueryType.Projects, QueryType.ProjectItems];

function hasRepo(qt?: string) {
  return !nonRepoTypes.includes(qt as QueryType);
}

export default QueryEditor;
