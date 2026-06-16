import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { RightColumnWidth } from './QueryEditor';
import { components } from '../components/selectors';
import type { RepositoryOptions } from '../types/query';
import { EditorField } from '@grafana/plugin-ui';

interface Props extends RepositoryOptions {
  onChange: (value: RepositoryOptions) => void;
}

export const QueryEditorRepository = (props: Props) => {
  const [repository, setRepository] = useState<string>(props.repository || '');
  // Track previous props to sync state during render (avoids extra render pass from useEffect)
  const [prevRepository, setPrevRepository] = useState(props.repository);

  if (props.repository !== prevRepository) {
    setPrevRepository(props.repository);
    setRepository(props.repository || '');
  }

  return (
    <EditorField label="Repository" tooltip="The name of the GitHub repository">
      <Input
        aria-label={components.QueryEditor.Repository.input}
        width={RightColumnWidth}
        value={repository}
        onChange={(el) => setRepository(el.currentTarget.value)}
        onBlur={(el) =>
          props.onChange({
            ...props,
            repository: el.currentTarget.value,
          })
        }
      />
    </EditorField>
  );
};

export const QueryEditorOwner = (props: Props) => {
  const [owner, setOwner] = useState<string>(props.owner || '');
  // Track previous props to sync state during render (avoids extra render pass from useEffect)
  const [prevOwner, setPrevOwner] = useState(props.owner);
  if (props.owner !== prevOwner) {
    setPrevOwner(props.owner);
    setOwner(props.owner || '');
  }
  return (
    <EditorField label="Owner" tooltip="The owner (organization or user) of the GitHub repository (example: 'grafana')">
      <Input
        aria-label={components.QueryEditor.Owner.input}
        width={RightColumnWidth}
        value={owner}
        onChange={(el) => setOwner(el.currentTarget.value)}
        onBlur={(el) =>
          props.onChange({
            ...props,
            owner: el.currentTarget.value,
          })
        }
      />
    </EditorField>
  );
};
