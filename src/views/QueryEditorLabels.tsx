import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import type { LabelsOptions } from '../types/query';

interface Props extends LabelsOptions {
  onChange: (value: LabelsOptions) => void;
}

export const QueryEditorLabels = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <EditorRow>
      <EditorField
        label="Query (optional)"
        tooltip={() => (
          <>
            For more information, visit&nbsp;
            <a
              href="https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests"
              target="_blank"
              rel="noreferrer"
            >
              https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests
            </a>
          </>
        )}
      >
        <Input
          width={RightColumnWidth * 2 + 2}
          value={query}
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, query: el.currentTarget.value })}
        />
      </EditorField>
    </EditorRow>
  );
};
