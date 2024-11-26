import React, { useState } from 'react';
import { Input, Select, InlineField } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { PullRequestTimeField } from '../constants';
import type { PullRequestsOptions } from '../types/query';

interface Props extends PullRequestsOptions {
  onChange: (value: PullRequestsOptions) => void;
}

const timeFieldOptions: Array<SelectableValue<PullRequestTimeField>> = Object.keys(PullRequestTimeField)
  .filter((_, i) => PullRequestTimeField[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${PullRequestTimeField[i]}`,
      value: i as PullRequestTimeField,
    };
  });

const defaultTimeField = timeFieldOptions[0].value;

const QueryEditorPullRequests = (props: Props) => {
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
          width={RightColumnWidth}
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
        tooltip="The time field to filter on the time range. WARNING: If selecting 'None', be mindful of the amount of data being queried. On larger repositories, querying all pull requests could easily cause rate limiting"
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
      </InlineField>
    </>
  );
};

export default QueryEditorPullRequests;
