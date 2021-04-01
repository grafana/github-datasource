import React, { useState } from 'react';
import { Input } from '@grafana/ui';

import { QueryInlineField } from '../components/Forms';
import { ContributorsOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends ContributorsOptions {
  onChange: (value: ContributorsOptions) => void;
}

const QueryEditorContributors = (props: Props) => {
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

export default QueryEditorContributors;
