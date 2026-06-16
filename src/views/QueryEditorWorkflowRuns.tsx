import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import type { WorkflowRunsOptions } from 'types/query';

interface Props extends WorkflowRunsOptions {
  onChange: (value: WorkflowRunsOptions) => void;
}

export const QueryEditorWorkflowRuns = (props: Props) => {
  const [workflow, setWorkflow] = useState<string | undefined>(props.workflow);
  const [branch, setBranch] = useState<string | undefined>(props.branch);

  return (
    <EditorRow>
      <EditorField label="Workflow" tooltip="The workflow id number or file name (e.g my-workflow.yml)">
        <Input
          value={workflow}
          width={RightColumnWidth}
          onChange={(el) => setWorkflow(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              workflow: el.currentTarget.value,
            })
          }
        />
      </EditorField>
      <EditorField label="Branch" tooltip="The branch to filter on (can be left empty)">
        <Input
          value={branch}
          width={RightColumnWidth}
          onChange={(el) => setBranch(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              branch: el.currentTarget.value,
            })
          }
        />
      </EditorField>
    </EditorRow>
  );
};
