import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import type { DeploymentsOptions } from 'types/query';

interface Props extends DeploymentsOptions {
  onChange: (value: DeploymentsOptions) => void;
}

export const QueryEditorDeployments = (props: Props) => {
  const { sha: initialSha, gitRef: initialRef, task: initialTask, environment: initialEnvironment } = props;
  const [sha, setSha] = useState<string | undefined>(initialSha);
  const [gitRef, setGitRef] = useState<string | undefined>(initialRef);
  const [task, setTask] = useState<string | undefined>(initialTask);
  const [environment, setEnvironment] = useState<string | undefined>(initialEnvironment);

  return (
    <EditorRow>
      <EditorField label="SHA" tooltip="Filter deployments by the SHA recorded at creation time (optional)">
        <Input
          value={sha}
          width={RightColumnWidth}
          onChange={(el) => setSha(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              sha: el.currentTarget.value || undefined,
            })
          }
        />
      </EditorField>
      <EditorField label="Ref" tooltip="Filter by ref name (branch, tag, or SHA) (optional)">
        <Input
          value={gitRef}
          width={RightColumnWidth}
          onChange={(el) => setGitRef(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              gitRef: el.currentTarget.value || undefined,
            })
          }
        />
      </EditorField>
      <EditorField label="Task" tooltip="Filter by task name (e.g., 'deploy', 'deploy:migrations') (optional)">
        <Input
          value={task}
          width={RightColumnWidth}
          onChange={(el) => setTask(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              task: el.currentTarget.value || undefined,
            })
          }
        />
      </EditorField>
      <EditorField
        label="Environment"
        tooltip="Filter by environment name (e.g., 'production', 'staging', 'qa') (optional)"
      >
        <Input
          value={environment}
          width={RightColumnWidth}
          onChange={(el) => setEnvironment(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              environment: el.currentTarget.value || undefined,
            })
          }
        />
      </EditorField>
    </EditorRow>
  );
};
