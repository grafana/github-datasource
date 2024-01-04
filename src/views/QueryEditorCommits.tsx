import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';

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
      <InlineField labelWidth={LeftColumnWidth * 2} label="Ref (Branch / Tag)">
        <Input
          aria-label={selectors.components.QueryEditor.Ref.input}
          width={RightColumnWidth}
          value={ref}
          onChange={(el) => setRef(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, gitRef: el.currentTarget.value })}
        />
      </InlineField>
    </>
  );
};

export default QueryEditorCommits;
