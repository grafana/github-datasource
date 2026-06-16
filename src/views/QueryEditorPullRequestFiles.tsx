import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import type { Pull_Request_FilesOption } from '../types/query';

interface Props extends Pull_Request_FilesOption {
  onChange: (value: Pull_Request_FilesOption) => void;
}

export const QueryEditorPullRequestFiles = (props: Props) => {
  const [prNumber, setPrNumber] = useState<string>(props.prNumber !== undefined ? String(props.prNumber) : '');
  return (
    <EditorRow>
      <EditorField label="Pull Request Number" tooltip="The pull request number to retrieve changed files for">
        <Input
          width={RightColumnWidth}
          type="number"
          value={prNumber}
          placeholder="e.g. 42"
          onChange={(el) => setPrNumber(el.currentTarget.value)}
          onBlur={(el) => {
            const parsed = parseInt(el.currentTarget.value, 10);
            props.onChange({ prNumber: isNaN(parsed) ? undefined : parsed });
          }}
        />
      </EditorField>
    </EditorRow>
  );
};

export default QueryEditorPullRequestFiles;
