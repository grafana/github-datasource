import React, { useState } from 'react';
import { Input } from '@grafana/ui';

import { QueryInlineField } from '../components/Forms';
import { LabelsOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends LabelsOptions {
  onChange: (value: LabelsOptions) => void;
}

const QueryEditorLabels = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <>
      <QueryInlineField labelWidth={LeftColumnWidth} label="Query (optional)">
        <Input
          css=""
          width={RightColumnWidth}
          value={query}
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, query: el.currentTarget.value })}
        />
      </QueryInlineField>
    </>
  );
};

export default QueryEditorLabels;
