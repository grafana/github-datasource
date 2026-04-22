import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { RightColumnWidth } from './QueryEditor';
import type { ContributorsOptions } from '../types/query';
import { EditorField, EditorRow } from '@grafana/plugin-ui';

interface Props extends ContributorsOptions {
  onChange: (value: ContributorsOptions) => void;
}

export const QueryEditorContributors = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <EditorRow>
      <EditorField label="Query (optional)">
        <Input
          width={RightColumnWidth}
          value={query}
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, query: el.currentTarget.value })}
        />
      </EditorField>
    </EditorRow>
  );
};
