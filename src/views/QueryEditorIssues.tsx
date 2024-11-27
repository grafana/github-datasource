import React, { useState } from 'react';
import { Input, Select, InlineField } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { components } from 'components/selectors';
import { IssueTimeField } from '../constants';
import type { IssuesOptions } from '../types/query';

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
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Query"
        tooltip={() => (
          <>
            For more information, visit&nbsp;
            <a
              href="https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests"
              target="_blank"
              rel="noreferrer"
            >
              https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests
            </a>
          </>
        )}
        interactive={true}
      >
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
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Time Field"
        tooltip="The time field to filter on the time range"
      >
        <Select
          data-testid={components.QueryEditor.Issues.timeFieldInput}
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
      </InlineField>
    </>
  );
};

export default QueryEditorIssues;
