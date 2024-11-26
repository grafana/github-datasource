import React, { useState, useEffect } from 'react';
import { Input, InlineLabel } from '@grafana/ui';
import { QueryEditorRow } from '../components/Forms';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { components } from '../components/selectors';
import type { RepositoryOptions } from '../types/query';

interface Props extends RepositoryOptions {
  onChange: (value: RepositoryOptions) => void;
}

const QueryEditorRepositories = (props: Props) => {
  const [repository, setRepository] = useState<string>(props.repository || '');
  const [owner, setOwner] = useState<string>(props.owner || '');

  useEffect(() => {
    setRepository(props.repository || '');
  }, [props.repository]);

  useEffect(() => {
    setOwner(props.owner || '');
  }, [props.owner]);

  return (
    <QueryEditorRow>
      <InlineLabel
        tooltip="The owner (organization or user) of the GitHub repository (example: 'grafana')"
        width={LeftColumnWidth * 2}
      >
        Owner
      </InlineLabel>
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
      <InlineLabel tooltip="The name of the GitHub repository" width={LeftColumnWidth * 2}>
        Repository
      </InlineLabel>
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
    </QueryEditorRow>
  );
};

export default QueryEditorRepositories;
