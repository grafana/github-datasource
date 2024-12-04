import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { WorkflowRunsOptions } from 'types/query';

interface Props extends WorkflowRunsOptions {
  onChange: (value: WorkflowRunsOptions) => void;
}

const QueryEditorWorkflowRuns = (props: Props) => {
  const [workflow, setWorkflow] = useState<string | undefined>(props.workflow);
  const [branch, setBranch] = useState<string | undefined>(props.branch);

  return (
    <>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
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
      </InlineField>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Branch"
        tooltip="The branch to filter on (can be left empty)"
      >
        <Input
          value={branch}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={(el) => setBranch(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              branch: el.currentTarget.value,
            })
          }
        />
      </InlineField>
    </>
  );
};

export default QueryEditorWorkflowRuns;
