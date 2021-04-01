import React, { useState } from 'react';

import { Input, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

import { QueryInlineField } from '../components/Forms';
import { IssuesOptions, IssueTimeField } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends IssuesOptions {
  onChange: (value: IssuesOptions) => void;
}

const timeFieldOptions: Array<SelectableValue<IssueTimeField>> = Object.keys(IssueTimeField)
  .filter((_, i) => IssueTimeField[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${IssueTimeField[i]}`,
      value: i as IssueTimeField,
    };
  });

const defaultTimeField = 0 as IssueTimeField;

const QueryEditorIssues = (props: Props) => {
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
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              query: el.currentTarget.value,
            })
          }
        />
      </QueryInlineField>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Time Field"
        tooltip="The time field to filter on th time range"
      >
        <Select
          width={RightColumnWidth}
          options={timeFieldOptions}
          value={props.timeField || defaultTimeField}
          onChange={(opt) =>
            props.onChange({
              ...props,
              timeField: opt.value,
            })
          }
        />
      </QueryInlineField>
    </>
  );
};

export default QueryEditorIssues;
