import React from 'react';

import { Input } from '@grafana/ui';

import { QueryInlineField } from '../components/Forms';
import { IssuesOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { DataSource } from '../DataSource';

interface Props extends IssuesOptions {
  owner: string;
  repository: string;
  datasource: DataSource;
  onChange: (value: IssuesOptions) => void;
}

export default (props: Props) => {
  return (
    <>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Query"
        tooltip="For more information, visit https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests"
      >
        <Input
          css=""
          value={props.query}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={el =>
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
