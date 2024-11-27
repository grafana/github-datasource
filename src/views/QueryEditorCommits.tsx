import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { components } from 'components/selectors';
import type { CommitsOptions } from '../types/query';

interface Props extends CommitsOptions {
  onChange: (value: CommitsOptions) => void;
}

const QueryEditorCommits = (props: Props) => {
  const [ref, setRef] = useState<string>(props.gitRef || '');
  return (
    <>
      <InlineField labelWidth={LeftColumnWidth * 2} label="Ref (Branch / Tag)">
        <Input
          aria-label={components.QueryEditor.Ref.input}
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
