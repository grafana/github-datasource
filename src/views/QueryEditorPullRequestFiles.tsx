import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { LeftColumnWidth, RightColumnWidth } from './QueryEditor';
import type { PullRequestFilesOptions } from '../types/query';

interface Props extends PullRequestFilesOptions {
  onChange: (value: PullRequestFilesOptions) => void;
}

const QueryEditorPullRequestFiles = (props: Props) => {
  const [prNumber, setPrNumber] = useState<string>(
    props.prNumber !== undefined ? String(props.prNumber) : ''
  );
  return (
    <>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Pull Request Number"
        tooltip="The pull request number to retrieve changed files for"
      >
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
      </InlineField>
    </>
  );
};

export default QueryEditorPullRequestFiles;
