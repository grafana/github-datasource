import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import type { WorkflowUsageOptions } from 'types/query';

interface Props extends WorkflowUsageOptions {
  onChange: (value: WorkflowUsageOptions) => void;
}

export const QueryEditorWorkflowUsage = (props: Props) => {
  const [workflow, setWorkflow] = useState<string | undefined>(props.workflow);

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
    </EditorRow>
  );
};
