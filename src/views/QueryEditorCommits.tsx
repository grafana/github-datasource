import React, { useState } from 'react';
import { Input } from '@grafana/ui';

import { QueryInlineField } from '../components/Forms';
import { CommitsOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { selectors } from 'components/selectors';

interface Props extends CommitsOptions {
  onChange: (value: CommitsOptions) => void;
}

const QueryEditorCommits = (props: Props) => {
  const [ref, setRef] = useState<string>(props.gitRef || '');
  return (
    <>
      <QueryInlineField labelWidth={LeftColumnWidth} label="Ref (Branch / Tag)">
        <Input
          aria-label={selectors.components.QueryEditor.Ref.input}
          css=""
          width={RightColumnWidth}
          value={ref}
          onChange={(el) => setRef(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, gitRef: el.currentTarget.value })}
        />
      </QueryInlineField>
    </>
  );
};

export default QueryEditorCommits;
