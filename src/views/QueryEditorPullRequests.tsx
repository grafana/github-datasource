import React from 'react';

import { Input, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

import QueryEditorRepository from './QueryEditorRepository';
import { QueryInlineField } from '../components/Forms';
import { PullRequestsOptions, PullRequestTimeField } from '../query';

import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends PullRequestsOptions {
  onChange: (value: PullRequestsOptions) => void;
};

const timeFieldOptions: Array<SelectableValue<PullRequestTimeField>> = Object.keys(PullRequestTimeField)
  .filter((_, i) => PullRequestTimeField[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${PullRequestTimeField[i]}`,
      value: i as PullRequestTimeField,
    };
  });

export default (props: Props) => {
  return (
    <>
      <QueryEditorRepository {...props} />
      <QueryInlineField labelWidth={LeftColumnWidth} label="Query" tooltip="For more information, visit https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests">
        <Input css='' value={props.query} width={RightColumnWidth} onChange={el => props.onChange({
          ...props,
          query: el.currentTarget.value,
        })}/>
      </QueryInlineField>
      <QueryInlineField labelWidth={LeftColumnWidth} label="Time Field" tooltip="The time field to filter on th time range">
        <Select
          width={RightColumnWidth}
          options={timeFieldOptions}
          value={props.timeField}
          onChange={opt => props.onChange({
            ...props,
            timeField: opt.value,
          })}
        />
      </QueryInlineField>
    </>
  );
}
