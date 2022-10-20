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
        css=""
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
    </QueryEditorRow>
  );
};

export default QueryEditorProjects;
