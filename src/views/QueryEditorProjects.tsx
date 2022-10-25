import React, { useState } from 'react';
import { Input, InlineFormLabel, RadioButtonGroup } from '@grafana/ui';

import { QueryEditorRow } from '../components/Forms';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { selectors } from '../components/selectors';
import { ProjectsOptions, ProjectQueryType } from 'types';

interface Props extends ProjectsOptions {
  onChange: (value: ProjectsOptions) => void;
}

const queryTypes = [
  { label: 'Organization', value: ProjectQueryType.ORG },
  { label: 'User', value: ProjectQueryType.USER},
];

const QueryEditorProjects = (props: Props) => {
  const [org, setOrg] = useState<string>(props.organization || '');
  const [user, setUser] = useState<string>(props.user || '');
  const [number, setNumber] = useState<number | undefined>(props.number);
  const [kind, setKind] = useState<ProjectQueryType>(props.kind || ProjectQueryType.ORG);
  const label = kind === ProjectQueryType.ORG ? 'Organization' : 'User';
  const tooltip = kind === ProjectQueryType.ORG ? 'The organization for the GitHub project (example: \'grafana\)' : 'The user who owns the Github project';

  return (
    <>
      <QueryEditorRow>
        <InlineFormLabel
          className="query-keyword"
          tooltip="The owner of the GitHub project"
          width={LeftColumnWidth}
        >
          Project Owner
        </InlineFormLabel>
        <div className="gf-form">
          <RadioButtonGroup<ProjectQueryType> options={queryTypes} value={kind} onChange={(v) => setKind(v!)} size={'md'} />
        </div>
      </QueryEditorRow>
      
      <QueryEditorRow>
        <InlineFormLabel
          className="query-keyword"
          tooltip={tooltip}
          width={LeftColumnWidth}
        >
          {label}
        </InlineFormLabel>
        {kind == ProjectQueryType.ORG &&
          <Input
            aria-label={selectors.components.QueryEditor.Owner.input}
            width={RightColumnWidth}
            value={org}
            onChange={(el) => setOrg(el.currentTarget.value)}
            onBlur={(el) =>
              props.onChange({
                ...props,
                organization: el.currentTarget.value,
                kind
              })
            }
          />
        }
        {kind == ProjectQueryType.USER &&
          <Input
            aria-label={selectors.components.QueryEditor.Owner.input}
            width={RightColumnWidth}
            value={user}
            onChange={(el) => setUser(el.currentTarget.value)}
            onBlur={(el) =>
              props.onChange({
                ...props,
                user: el.currentTarget.value,
                kind
              })
            }
          />
        }
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
              number: num(el.currentTarget.value)
            })
          }
        />
      </QueryEditorRow>
    </>
  );
};

function num(v: string) {
  const val = parseInt(v);
  return isNaN(val) ? undefined : val;
}

export default QueryEditorProjects;
