import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import type { Pull_Requests_For_BranchesOptions } from '../types/query';

interface Props extends Pull_Requests_For_BranchesOptions {
  onChange: (value: Pull_Requests_For_BranchesOptions) => void;
}

export const QueryEditorPullRequestsForBranches = ({ branches = [], onChange }: Props) => {
  const [value, setValue] = useState<string>(branches.join(', '));

  return (
    <EditorRow>
      <EditorField label="Branches" tooltip="Comma-separated head branch names to look up, at most 100">
        <Input
          aria-label="Branches"
          width={RightColumnWidth * 2}
          value={value}
          placeholder="e.g. feat/new-panel, fix/typo"
          onChange={(e) => setValue(e.currentTarget.value)}
          onBlur={() =>
            onChange({
              branches: value
                .split(',')
                .map((branch) => branch.trim())
                .filter(Boolean),
            })
          }
        />
      </EditorField>
    </EditorRow>
  );
};
