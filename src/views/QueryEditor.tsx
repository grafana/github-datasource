import React, { ReactNode, useCallback } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { Combobox, ComboboxOption } from '@grafana/ui';
import { EditorField, EditorRow, EditorRows } from '@grafana/plugin-ui';

import { GitHubDataSource } from '../DataSource';
import { isValid } from '../validation';
import { components } from '../components/selectors';

import { QueryEditorOwner, QueryEditorRepository } from './QueryEditorRepository';
import { QueryEditorCommits } from './QueryEditorCommits';
import { QueryEditorIssues } from './QueryEditorIssues';
import { QueryEditorMilestones } from './QueryEditorMilestones';
import { QueryEditorPullRequests } from './QueryEditorPullRequests';
import { QueryEditorPullRequestReviews } from './QueryEditorPullRequestReviews';
import { QueryEditorContributors } from './QueryEditorContributors';
import { QueryEditorLabels } from './QueryEditorLabels';
import { QueryEditorPackages } from './QueryEditorPackages';
import { QueryEditorProjects } from './QueryEditorProjects';
import { QueryEditorWorkflows } from './QueryEditorWorkflows';
import { QueryEditorWorkflowUsage } from './QueryEditorWorkflowUsage';
import { QueryEditorWorkflowRuns } from './QueryEditorWorkflowRuns';
import { QueryEditorCodeScanning } from './QueryEditorCodeScanning';
import { QueryEditorDeployments } from './QueryEditorDeployments';

import { DefaultQueryType, QueryTypes } from '../constants';

import type { QueryType, GitHubQuery } from '../types/query';
import type { GitHubDataSourceOptions } from '../types/config';

interface Props extends QueryEditorProps<GitHubDataSource, GitHubQuery, GitHubDataSourceOptions> {
  queryTypes?: QueryType[];
}
export const LeftColumnWidth = 10;
export const RightColumnWidth = 36;

const queryEditors: Record<QueryType, { component: (props: Props, onChange: (val: any) => void) => ReactNode }> = {
  ['Repositories']: { component: () => <></> },
  ['GraphQL']: { component: () => <></> },
  ['Organizations']: { component: () => <></> },
  ['ProjectItems']: { component: () => <></> },
  ['Tags']: { component: () => <></> },
  ['Releases']: { component: () => <></> },
  ['Vulnerabilities']: { component: () => <></> },
  ['Stargazers']: { component: () => <></> },
  ['Labels']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorLabels {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Contributors']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorContributors {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Code_Scanning']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorCodeScanning {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Commits']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorCommits {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Milestones']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorMilestones {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Issues']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorIssues {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Packages']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorPackages {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Pull_Requests']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorPullRequests {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Pull_Request_Reviews']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorPullRequestReviews {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Projects']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorProjects {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Workflows']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorWorkflows {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Workflow_Usage']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorWorkflowUsage {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Workflow_Runs']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorWorkflowRuns {...(props.query.options || {})} onChange={onChange} />
    ),
  },
  ['Deployments']: {
    component: (props: Props, onChange: (val: any) => void) => (
      <QueryEditorDeployments {...(props.query.options || {})} onChange={onChange} />
    ),
  },
};

const queryTypeOptions: Array<SelectableValue<QueryType>> = QueryTypes.map((v) => {
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
      <EditorRows>
        <EditorRow>
          <EditorField label="Query Type" tooltip={'What resource are you querying for?'}>
            <div aria-label={components.QueryEditor.QueryType.container.ariaLabel}>
              <Combobox<QueryType>
                width={RightColumnWidth}
                options={
                  queryTypeOptions.filter(
                    (v) =>
                      queryTypes.includes(v.value!) &&
                      v.value !== 'Organizations' &&
                      v.value !== 'GraphQL' &&
                      v.value !== 'ProjectItems'
                  ) as Array<ComboboxOption<QueryType>>
                }
                value={props.query.queryType}
                onChange={(val) => onKeyChange('queryType', val.value || DefaultQueryType)}
              />
            </div>
          </EditorField>
          {hasRepo(props.query.queryType) && (
            <QueryEditorOwner
              owner={props.query.owner}
              onChange={(repo) => {
                onChange({
                  ...props.query,
                  ...repo,
                });
              }}
            />
          )}
          {hasRepo(props.query.queryType) && (
            <QueryEditorRepository
              repository={props.query.repository}
              onChange={(repo) => {
                onChange({
                  ...props.query,
                  ...repo,
                });
              }}
            />
          )}
        </EditorRow>
        {queryEditor ? (
          queryEditor.component(props, (value: any) => onKeyChange('options', !!value ? value : undefined))
        ) : (
          <span>Unsupported Query Type</span>
        )}
      </EditorRows>
    </>
  );
};

const nonRepoTypes = ['Projects', 'ProjectItems'];

function hasRepo(qt?: string) {
  return !nonRepoTypes.includes(qt as QueryType);
}

export default QueryEditor;
