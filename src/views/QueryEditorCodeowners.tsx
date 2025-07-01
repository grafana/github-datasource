import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { CodeownersOptions } from '../types/query';

interface Props extends CodeownersOptions {
  onChange: (value: CodeownersOptions) => void;
}

const QueryEditorCodeowners = (props: Props) => {
  const [filePath, setFilePath] = useState<string>(props.filePath || '');
  
  return (
    <>
      <InlineField 
        labelWidth={LeftColumnWidth * 2} 
        label="File Path (optional)"
        tooltip="Optional file path to find owners for (e.g., 'src/main.go', 'docs/README.md'). Leave empty to show all CODEOWNERS entries."
      >
        <Input
          aria-label="File Path"
          width={RightColumnWidth}
          value={filePath}
          placeholder="e.g., src/main.go or docs/*.md"
          onChange={(el) => setFilePath(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, filePath: el.currentTarget.value })}
        />
      </InlineField>
    </>
  );
};

export default QueryEditorCodeowners; 