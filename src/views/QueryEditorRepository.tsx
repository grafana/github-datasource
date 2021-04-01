import React, { useState } from 'react';
import { Input, InlineFormLabel } from '@grafana/ui';

import { QueryEditorRow } from '../components/Forms';
import { RepositoryOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { selectors } from '../components/selectors';

interface Props extends RepositoryOptions {
  onChange: (value: RepositoryOptions) => void;
}

const QueryEditorRepositories = (props: Props) => {
  const [repository, setRepository] = useState<string>(props.repository || '');
  const [owner, setOwner] = useState<string>(props.owner || '');

  return (
    <QueryEditorRow>
      <InlineFormLabel
        className="query-keyword"
        tooltip="The owner (organization or user) of the GitHub repository (example: 'grafana')"
        width={LeftColumnWidth}
      >
        Owner
      </InlineFormLabel>
      <Input
        aria-label={selectors.components.QueryEditor.Owner.input}
        css=""
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
      <InlineFormLabel className="query-keyword" tooltip="The name of the GitHub repository" width={LeftColumnWidth}>
        Repository
      </InlineFormLabel>
      <Input
        aria-label={selectors.components.QueryEditor.Repository.input}
        css=""
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
