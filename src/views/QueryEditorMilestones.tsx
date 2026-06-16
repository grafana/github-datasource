import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import type { MilestonesOptions } from '../types/query';

interface Props extends MilestonesOptions {
  onChange: (value: MilestonesOptions) => void;
}

export const QueryEditorMilestones = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <EditorRow>
      <EditorField label="Query" tooltip="Query milestones by title">
        <Input
          value={query}
          width={RightColumnWidth * 2 + 2}
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              query: el.currentTarget.value,
            })
          }
        />
      </EditorField>
    </EditorRow>
  );
};
