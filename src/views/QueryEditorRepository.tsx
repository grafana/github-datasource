import React from 'react';
import { Input, InlineFormLabel } from '@grafana/ui';

import { QueryEditorRow, QueryRowTerminator } from '../components/Forms';
import { RepositoryOptions } from '../query';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends RepositoryOptions {
  onChange: (value: RepositoryOptions) => void;
};

export default (props: Props) => {
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
        css=""
        width={RightColumnWidth}
        value={props.owner}
        onChange={el =>
          props.onChange({
            ...props,
            owner: el.currentTarget.value,
          })
        }
      />
      <InlineFormLabel
        className="query-keyword"
        tooltip="The name of the GitHub repository"
        width={LeftColumnWidth}>
        Repository 
      </InlineFormLabel>
      <Input
        css=""
        width={RightColumnWidth}
        value={props.repository}
        onChange={el =>
          props.onChange({
            ...props,
            repository: el.currentTarget.value,
          })
        }
      />
      <QueryRowTerminator />
    </QueryEditorRow>
  );
}
