import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { MilestonesOptions } from '../types/query';

interface Props extends MilestonesOptions {
  onChange: (value: MilestonesOptions) => void;
}

const QueryEditorMilestones = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <>
      <InlineField labelWidth={LeftColumnWidth * 2} label="Query" tooltip="Query milestones by title">
        <Input
          value={query}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              query: el.currentTarget.value,
            })
          }
        />
      </InlineField>
    </>
  );
};

export default QueryEditorMilestones;
