import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { BranchesOptions } from '../types/query';
import { LeftColumnWidth, RightColumnWidth } from './QueryEditor';

interface Props extends BranchesOptions {
  onChange: (value: BranchesOptions) => void;
}

export const QueryEditorBranches = ({ query = '', onChange }: Props) => {
  const [filter, setFilter] = useState<string>(query);

  return (
    <EditorRow>
      <EditorField
        label="Filter"
        tooltip="Filter branches by name prefix (e.g. release/ matches all release/* branches). Leave empty to list all branches."
        width={LeftColumnWidth}
      >
        <Input
          aria-label="Branch filter"
          placeholder="release/"
          value={filter}
          onChange={(e) => setFilter(e.currentTarget.value)}
          onBlur={() => onChange({ query: filter })}
          width={RightColumnWidth}
        />
      </EditorField>
    </EditorRow>
  );
};
