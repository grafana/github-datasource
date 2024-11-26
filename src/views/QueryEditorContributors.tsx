import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { ContributorsOptions } from '../types/query';

interface Props extends ContributorsOptions {
  onChange: (value: ContributorsOptions) => void;
}

const QueryEditorContributors = (props: Props) => {
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

export default QueryEditorContributors;
