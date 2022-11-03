import React, { useState } from 'react';
import { Input, InlineFormLabel, RadioButtonGroup } from '@grafana/ui';

import { QueryEditorRow } from '../components/Forms';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { selectors } from '../components/selectors';
import { ProjectsOptions, ProjectQueryType } from 'types';
import { SelectableValue } from '@grafana/data';
import { Filter, Filters } from 'components/Filters';

interface Props extends ProjectsOptions {
  onChange: (value: ProjectsOptions) => void;
}

const queryTypes = [
  { label: 'Organization', value: ProjectQueryType.ORG },
  { label: 'User', value: ProjectQueryType.USER },
];

const filters: Array<SelectableValue<string>> = [
  { label: 'Type', value: 'type' },
  { label: 'Created At', value: 'created_at' },
  { label: 'Title', value: 'Title' },
  { label: 'Assignees', value: 'Assignees' },
  { label: 'Status', value: 'Status' },
  { label: 'Labels', value: 'Labels' },
  { label: 'Reviewers', value: 'Reviewers' },
  { label: 'Milestone', value: 'Milestone' },
  { label: 'Iteration', value: 'Iteration' },
];

const ops: Array<SelectableValue<string>> = [
  { label: 'Equal', value: '=' },
  { label: 'Not Equal', value: '!=' },
  { label: 'Greater Than', value: '>' },
  { label: 'Less Than', value: '<' },
  { label: 'Less Than or Equal', value: '<=' },
  { label: 'Greater Than or Equal', value: '>=' },
  { label: 'Contains', value: '~' },
];

const fetchFilters = async (key?: string) => (key ? [] : filters);

const QueryEditorProjects = (props: Props) => {
  const [org, setOrg] = useState<string>(props.organization || '');
  const [user, setUser] = useState<string>(props.user || '');
  const [number, setNumber] = useState<number | string | undefined>(props.number);
  const [kind, setKind] = useState<ProjectQueryType>(props.kind || ProjectQueryType.ORG);
  const [filters, setFilters] = useState<Filter[]>(props.filters || []);
  const label = kind === ProjectQueryType.ORG ? 'Organization' : 'User';
  const tooltip =
    kind === ProjectQueryType.ORG
      ? "The organization for the GitHub project (example: 'grafana)"
      : 'The user who owns the Github project';

  return (
    <>
      <QueryEditorRow>
        <InlineFormLabel className="query-keyword" tooltip="The owner of the GitHub project" width={LeftColumnWidth}>
          Project Owner
        </InlineFormLabel>
        <div className="gf-form">
          <RadioButtonGroup<ProjectQueryType>
            options={queryTypes}
            value={kind}
            onChange={(v) => setKind(v!)}
            size={'md'}
          />
        </div>
      </QueryEditorRow>

      <QueryEditorRow>
        <InlineFormLabel className="query-keyword" tooltip={tooltip} width={LeftColumnWidth}>
          {label}
        </InlineFormLabel>
        {kind === ProjectQueryType.ORG && (
          <Input
            aria-label={selectors.components.QueryEditor.Owner.input}
            width={RightColumnWidth}
            value={org}
            onChange={(el) => setOrg(el.currentTarget.value)}
            onBlur={(el) =>
              props.onChange({
                ...props,
                organization: el.currentTarget.value,
                kind,
              })
            }
          />
        )}
        {kind === ProjectQueryType.USER && (
          <Input
            aria-label={selectors.components.QueryEditor.Owner.input}
            width={RightColumnWidth}
            value={user}
            onChange={(el) => setUser(el.currentTarget.value)}
            onBlur={(el) =>
              props.onChange({
                ...props,
                user: el.currentTarget.value,
                kind,
              })
            }
          />
        )}
      </QueryEditorRow>

      <QueryEditorRow>
        <InlineFormLabel
          className="query-keyword"
          tooltip="The project number for the GitHub project (example: 123)"
          width={LeftColumnWidth}
        >
          Project Number
        </InlineFormLabel>
        <Input
          aria-label={selectors.components.QueryEditor.Number.input}
          width={RightColumnWidth}
          value={number}
          onChange={(el) => setNumber(num(el.currentTarget.value))}
          onBlur={(el) =>
            props.onChange({
              ...props,
              number: num(el.currentTarget.value),
            })
          }
        />
      </QueryEditorRow>

      <QueryEditorRow>
        <div className="gf-form">
          <InlineFormLabel className="query-keyword" width={LeftColumnWidth}>
            Filters
          </InlineFormLabel>
          <Filters
            onChange={(filters: Filter[]) => {
              setFilters(filters);
              props.onChange({
                ...props,
                filters,
              });
            }}
            loadOptions={fetchFilters}
            value={filters}
            ops={ops}
          ></Filters>
        </div>
      </QueryEditorRow>
    </>
  );
};

function num(v: string) {
  if (v.includes('$')) {
    return v;
  }
  const val = parseInt(v, 10);
  return isNaN(val) ? undefined : val;
}

export default QueryEditorProjects;
