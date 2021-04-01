import React, { useState } from 'react';

import { Input } from '@grafana/ui';

import { QueryInlineField } from '../components/Forms';
import { MilestonesOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends MilestonesOptions {
  onChange: (value: MilestonesOptions) => void;
}

const QueryEditorMilestones = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <>
      <QueryInlineField labelWidth={LeftColumnWidth} label="Query" tooltip="Query milestones by title">
        <Input
          css=""
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
      </QueryInlineField>
    </>
  );
};

export default QueryEditorMilestones;
