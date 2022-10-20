import React, { useState } from 'react';
import { Input, InlineFormLabel } from '@grafana/ui';

import { QueryEditorRow } from '../components/Forms';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { selectors } from '../components/selectors';
import { ProjectsOptions } from 'types';

interface Props extends ProjectsOptions {
  onChange: (value: ProjectsOptions) => void;
}

const QueryEditorProjects = (props: Props) => {
  const [org, setOrg] = useState<string>(props.organization || '');
  const [number, setNumber] = useState<number | undefined>(props.number);

  return (
    <QueryEditorRow>
      <InlineFormLabel
        className="query-keyword"
        tooltip="The organization for the GitHub project (example: 'grafana')"
        width={LeftColumnWidth}
      >
        Organization
      </InlineFormLabel>
      <Input
        aria-label={selectors.components.QueryEditor.Owner.input}
        width={RightColumnWidth}
        value={org}
        onChange={(el) => setOrg(el.currentTarget.value)}
        onBlur={(el) =>
          props.onChange({
            ...props,
            organization: el.currentTarget.value,
          })
        }
      />
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
  );
};

function num(v: string) {
  const val = parseInt(v);
  return isNaN(val) ? undefined : val;
}

export default QueryEditorProjects;
