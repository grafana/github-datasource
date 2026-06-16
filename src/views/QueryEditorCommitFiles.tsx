import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import type { CommitFilesOptions } from '../types/query';

interface Props extends CommitFilesOptions {
  onChange: (value: CommitFilesOptions) => void;
}

export const QueryEditorCommitFiles = (props: Props) => {
  const [commitSha, setCommitSha] = useState<string>(props.commitSha || '');
  return (
    <EditorRow>
      <EditorField label="Commit SHA" tooltip="The commit SHA to retrieve changed files for">
        <Input
          width={RightColumnWidth}
          value={commitSha}
          placeholder="e.g. abc123def456"
          onChange={(el) => setCommitSha(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ commitSha: el.currentTarget.value })}
        />
      </EditorField>
    </EditorRow>
  );
};
