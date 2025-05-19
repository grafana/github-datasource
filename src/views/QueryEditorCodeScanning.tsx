import React, { useState } from 'react';
import { LeftColumnWidth, RightColumnWidth } from './QueryEditor';
import { InlineField, Input } from '@grafana/ui';
import { components } from '../components/selectors';
import { CodeScanningOptions } from '../types/query';

interface Props extends CodeScanningOptions {
  onChange: (value: CodeScanningOptions) => void;
}

const QueryEditorCodeScanning = (props: Props) => {
  const [state, setState] = useState<string>(props.state || 'open');
  const [gitRef, setGitRef] = useState<string>(props.gitRef || '');
  return (
    <>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="State"
        tooltip="Can be one of: open, closed, dismissed, fixed. Default: open"
      >
        <Input
          aria-label={components.QueryEditor.CodeScanState.input}
          width={RightColumnWidth}
          value={state}
          onChange={(el) => setState(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, gitRef, state: el.currentTarget.value })}
        />
      </InlineField>
      <InlineField labelWidth={LeftColumnWidth * 2} label="Ref (Branch / Tag)">
        <Input
          aria-label={components.QueryEditor.Ref.input}
          width={RightColumnWidth}
          value={gitRef}
          onChange={(el) => setGitRef(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, state, gitRef: el.currentTarget.value })}
        />
      </InlineField>
    </>
  );
};
export default QueryEditorCodeScanning;
