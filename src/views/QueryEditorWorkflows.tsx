import React from 'react';
import { Select, InlineField } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { WorkflowsTimeField } from './../constants';
import type { WorkflowsOptions } from 'types/query';

interface Props extends WorkflowsOptions {
  onChange: (value: WorkflowsOptions) => void;
}

const timeFieldOptions: Array<SelectableValue<WorkflowsTimeField>> = Object.keys(WorkflowsTimeField)
  .filter((_, i) => WorkflowsTimeField[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${WorkflowsTimeField[i]}`,
      value: i as WorkflowsTimeField,
    };
  });

const defaultTimeField = 0 as WorkflowsTimeField;

const QueryEditorWorkflows = (props: Props) => {
  return (
    <>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Time Field"
        tooltip="The time field to filter on the time range"
      >
        <Select
          width={RightColumnWidth}
          options={timeFieldOptions}
          value={props.timeField || defaultTimeField}
          onChange={(opt) => {
            props.onChange({
              ...props,
              timeField: opt.value,
            });
          }}
        />
      </InlineField>
    </>
  );
};

export default QueryEditorWorkflows;
