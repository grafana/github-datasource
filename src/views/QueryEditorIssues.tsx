import React, { useState } from 'react';

import { Input } from '@grafana/ui';

import { QueryInlineField } from '../components/Forms';
import { IssuesOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends IssuesOptions {
  onChange: (value: IssuesOptions) => void;
}

export default (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Query"
        tooltip="For more information, visit https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests"
      >
        <Input
          css=""
          value={query}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={el => setQuery(el.currentTarget.value)}
          onBlur={el =>
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
