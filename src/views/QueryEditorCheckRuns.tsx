import React, { useState } from 'react';
import { Combobox, type ComboboxOption, Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import { components } from '../components/selectors';
import type { CheckRunsOptions } from '../types/query';

interface Props extends CheckRunsOptions {
  onChange: (value: CheckRunsOptions) => void;
}

const statusOptions: Array<ComboboxOption<string>> = [
  { label: 'Any', value: '' },
  { label: 'Queued', value: 'queued' },
  { label: 'In progress', value: 'in_progress' },
  { label: 'Completed', value: 'completed' },
];

const filterOptions: Array<ComboboxOption<string>> = [
  { label: 'Latest', value: 'latest' },
  { label: 'All', value: 'all' },
];

export const QueryEditorCheckRuns = (props: Props) => {
  const [gitRef, setGitRef] = useState<string>(props.gitRef || '');
  const [checkName, setCheckName] = useState<string>(props.checkName || '');

  return (
    <EditorRow>
      <EditorField label="Ref" tooltip="The commit SHA, branch (heads/<branch>), or tag (tags/<tag>) to list check runs for">
        <Input
          aria-label={components.QueryEditor.Ref.input}
          width={RightColumnWidth}
          value={gitRef}
          onChange={(el) => setGitRef(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, checkName, gitRef: el.currentTarget.value })}
        />
      </EditorField>
      <EditorField label="Check name" tooltip="Only return check runs with this exact name (can be left empty)">
        <Input
          width={RightColumnWidth}
          value={checkName}
          onChange={(el) => setCheckName(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, gitRef, checkName: el.currentTarget.value })}
        />
      </EditorField>
      <EditorField label="Status" tooltip="Only return check runs with this status">
        <Combobox<string>
          width={RightColumnWidth}
          options={statusOptions}
          value={props.status || ''}
          onChange={(value) => props.onChange({ ...props, gitRef, checkName, status: value?.value })}
        />
      </EditorField>
      <EditorField
        label="Filter"
        tooltip="Latest returns the most recent check run for each name, all returns every check run"
      >
        <Combobox<string>
          width={RightColumnWidth}
          options={filterOptions}
          value={props.filter || 'latest'}
          onChange={(value) => props.onChange({ ...props, gitRef, checkName, filter: value?.value })}
        />
      </EditorField>
    </EditorRow>
  );
};
