import React, { useState } from 'react';
import { Input, InlineLabel, RadioButtonGroup, InlineField } from '@grafana/ui';
import { QueryEditorRow } from '../components/Forms';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { components } from '../components/selectors';
import { SelectableValue } from '@grafana/data';
import { Filter, Filters } from 'components/Filters';
import { ProjectQueryType } from './../constants';
import type { ProjectsOptions } from 'types/query';

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
      ? "The organization for the GitHub project (example: 'grafana')"
      : 'The user who owns the GitHub project';

  return (
    <>
      <QueryEditorRow>
        <InlineLabel
          tooltip="The owner of the GitHub project."
          width={LeftColumnWidth * 2}
        >
          Project Owner
        </InlineLabel>
        <InlineField>
          <RadioButtonGroup<ProjectQueryType>
            options={queryTypes}
            value={kind}
            onChange={(v) => {
              setKind(v);
              props.onChange({
                ...props,
                kind: v,
              });
            }}
            size={'md'}
          />
        </InlineField>
      </QueryEditorRow>

      <QueryEditorRow>
        <InlineLabel tooltip={tooltip} width={LeftColumnWidth * 2}>
          {label}
        </InlineLabel>
        {kind === ProjectQueryType.ORG && (
          <Input
            aria-label={components.QueryEditor.Owner.input}
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
            aria-label={components.QueryEditor.Owner.input}
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
        <InlineLabel
          tooltip="The project number for the GitHub project (example: 123)."
          width={LeftColumnWidth * 2}
        >
          Project Number
        </InlineLabel>
        <Input
          aria-label={components.QueryEditor.Number.input}
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

      {/* Filters currently only apply to Project Items */}
      {number && (
        <QueryEditorRow>
          <InlineLabel width={LeftColumnWidth * 2}>Filters</InlineLabel>
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
        </QueryEditorRow>
      )}
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
