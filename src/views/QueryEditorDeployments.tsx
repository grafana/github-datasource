import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { DeploymentsOptions } from 'types/query';

interface Props extends DeploymentsOptions {
  onChange: (value: DeploymentsOptions) => void;
}

const QueryEditorDeployments = (props: Props) => {
  const { sha: initialSha, ref: initialRef, task: initialTask, environment: initialEnvironment } = props;
  const [sha, setSha] = useState<string | undefined>(initialSha);
  // eslint-disable-next-line react-hooks/refs -- 'ref' is a prop name from DeploymentsOptions, not a React ref
  const [gitRef, setGitRef] = useState<string | undefined>(initialRef);
  const [task, setTask] = useState<string | undefined>(initialTask);
  const [environment, setEnvironment] = useState<string | undefined>(initialEnvironment);

  return (
    <>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="SHA"
        tooltip="Filter deployments by the SHA recorded at creation time (optional)"
      >
        <Input
          value={sha}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={(el) => setSha(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              sha: el.currentTarget.value || undefined,
            })
          }
        />
      </InlineField>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Ref"
        tooltip="Filter by ref name (branch, tag, or SHA) (optional)"
      >
        <Input
          value={gitRef}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={(el) => setGitRef(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              ref: el.currentTarget.value || undefined,
            })
          }
        />
      </InlineField>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Task"
        tooltip="Filter by task name (e.g., 'deploy', 'deploy:migrations') (optional)"
      >
        <Input
          value={task}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={(el) => setTask(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              task: el.currentTarget.value || undefined,
            })
          }
        />
      </InlineField>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Environment"
        tooltip="Filter by environment name (e.g., 'production', 'staging', 'qa') (optional)"
      >
        <Input
          value={environment}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={(el) => setEnvironment(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              environment: el.currentTarget.value || undefined,
            })
          }
        />
      </InlineField>
    </>
  );
};

export default QueryEditorDeployments;
