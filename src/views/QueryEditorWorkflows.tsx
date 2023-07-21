import React from 'react';
import { QueryInlineField } from '../components/Forms';
import { Select } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { WorkflowsOptions, WorkflowsTimeField } from 'types';
import { SelectableValue } from '@grafana/data';

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
      <QueryInlineField
        labelWidth={LeftColumnWidth}
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
      </QueryInlineField>
    </>
  );
};

export default QueryEditorWorkflows;
