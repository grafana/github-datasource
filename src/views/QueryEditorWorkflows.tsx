import React from 'react';
import { Combobox, ComboboxOption } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import { WorkflowsTimeField } from './../constants';
import type { WorkflowsOptions } from 'types/query';

interface Props extends WorkflowsOptions {
  onChange: (value: WorkflowsOptions) => void;
}

const timeFieldOptions: Array<ComboboxOption<WorkflowsTimeField>> = Object.keys(WorkflowsTimeField)
  .filter((_, i) => WorkflowsTimeField[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${WorkflowsTimeField[i]}`,
      value: i as WorkflowsTimeField,
    };
  });

const defaultTimeField = WorkflowsTimeField.None;

export const QueryEditorWorkflows = (props: Props) => {
  return (
    <EditorRow>
      <EditorField
        label="Time Field"
        tooltip="Select 'None' to return all workflows, or choose a time field to filter by the dashboard time range"
      >
        <Combobox
          width={RightColumnWidth}
          options={timeFieldOptions}
          value={props.timeField !== undefined ? props.timeField : defaultTimeField}
          onChange={(opt) => {
            props.onChange({
              ...props,
              timeField: opt.value,
            });
          }}
        />
      </EditorField>
    </EditorRow>
  );
};
