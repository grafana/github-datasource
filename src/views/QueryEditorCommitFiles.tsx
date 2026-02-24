import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { LeftColumnWidth, RightColumnWidth } from './QueryEditor';
import type { CommitFilesOptions } from '../types/query';

interface Props extends CommitFilesOptions {
  onChange: (value: CommitFilesOptions) => void;
}

const QueryEditorCommitFiles = (props: Props) => {
  const [commitSha, setCommitSha] = useState<string>(props.commitSha || '');
  return (
    <>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Commit SHA"
        tooltip="The commit SHA to retrieve changed files for"
      >
        <Input
          width={RightColumnWidth}
          value={commitSha}
          placeholder="e.g. abc123def456"
          onChange={(el) => setCommitSha(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, commitSha: el.currentTarget.value })}
        />
      </InlineField>
    </>
  );
};

export default QueryEditorCommitFiles;
