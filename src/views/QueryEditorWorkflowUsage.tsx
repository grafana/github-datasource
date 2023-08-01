import React, { useState } from 'react';
import { QueryInlineField } from '../components/Forms';
import { Input } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { WorkflowUsageOptions } from 'types';

interface Props extends WorkflowUsageOptions {
  onChange: (value: WorkflowUsageOptions) => void;
}

const QueryEditorWorkflowUsage = (props: Props) => {
  const [workflow, setWorkflow] = useState<string | undefined>(props.workflow);

  return (
    <>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Workflow"
        tooltip="The workflow id number or file name (e.g my-workflow.yml)"
      >
        <Input
          value={workflow}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={(el) => setWorkflow(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              workflow: el.currentTarget.value,
            })
          }
        />
      </QueryInlineField>
    </>
  );
};

export default QueryEditorWorkflowUsage;
