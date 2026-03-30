import React, { useState } from 'react';
import { Input, InlineSwitch } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
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
          value={ref}
          onChange={(el) => setRef(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, gitRef: el.currentTarget.value })}
        />
      </EditorField>
      <EditorField
        label="Include File Changes"
        tooltip="Returns one row per changed file per commit. Makes one additional API call per commit — avoid large time ranges."
      >
        <InlineSwitch
          value={props.includeFiles || false}
          onChange={(el) => props.onChange({ gitRef: props.gitRef, includeFiles: el.currentTarget.checked })}
        />
      </EditorField>
    </EditorRow>
  );
};
