import React, { useState } from 'react';

import { Input, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

import { QueryInlineField } from '../components/Forms';
import { PullRequestsOptions, PullRequestTimeField } from '../types';

import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

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
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Query"
        tooltip="For more information, visit https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests"
      >
        <Input
          css=""
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
      </QueryInlineField>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
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
      </QueryInlineField>
    </>
  );
};

export default QueryEditorPullRequests;
