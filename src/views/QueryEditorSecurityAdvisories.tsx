import React, { useId } from 'react';
import { Combobox, ComboboxOption } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';

import { RightColumnWidth } from './QueryEditor';
import type { SecurityAdvisoriesOptions, SecurityAdvisoryState } from '../types/query';

interface Props extends SecurityAdvisoriesOptions {
  onChange: (options: SecurityAdvisoriesOptions) => void;
}

type StateChoice = SecurityAdvisoryState | 'all';

const stateOptions: Array<ComboboxOption<StateChoice>> = [
  { label: 'All', value: 'all' },
  { label: 'Triage', value: 'triage' },
  { label: 'Draft', value: 'draft' },
  { label: 'Published', value: 'published' },
  { label: 'Closed', value: 'closed' },
];

export function QueryEditorSecurityAdvisories({ state, onChange }: Props) {
  const stateId = useId();

  return (
    <EditorRow>
      <EditorField label="State" htmlFor={stateId} tooltip="Leave Repository empty to query the organization named in Owner.">
        <Combobox<StateChoice>
          id={stateId}
          options={stateOptions}
          value={state ?? 'all'}
          width={RightColumnWidth}
          onChange={({ value }) => onChange({ state: value === 'all' ? undefined : value })}
        />
      </EditorField>
    </EditorRow>
  );
}
