import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import { components } from 'components/selectors';
import type { CommitsOptions } from '../types/query';

interface Props extends CommitsOptions {
  onChange: (value: CommitsOptions) => void;
}

export const QueryEditorCommits = (props: Props) => {
  const [ref, setRef] = useState<string>(props.gitRef || '');
  return (
    <EditorRow>
      <EditorField label="Ref (Branch / Tag)">
        <Input
          aria-label={components.QueryEditor.Ref.input}
          width={RightColumnWidth}
          value={ref}
          onChange={(el) => setRef(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, gitRef: el.currentTarget.value })}
        />
      </EditorField>
    </EditorRow>
  );
};
