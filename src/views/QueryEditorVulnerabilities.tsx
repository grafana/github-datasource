import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';

import { LabelsOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends LabelsOptions {
  onChange: (value: LabelsOptions) => void;
}

const QueryEditorVulnerabilities = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <>
      <InlineField labelWidth={LeftColumnWidth * 2} label="Query (optional)">
        <Input
          width={RightColumnWidth}
          value={query}
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, query: el.currentTarget.value })}
        />
      </InlineField>
    </>
  );
};

export default QueryEditorVulnerabilities;
