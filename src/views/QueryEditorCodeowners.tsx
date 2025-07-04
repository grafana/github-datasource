import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { CodeownersOptions } from '../types/query';

interface Props extends CodeownersOptions {
  onChange: (value: CodeownersOptions) => void;
}

const QueryEditorCodeowners = (props: Props) => {
  const [filePath, setFilePath] = useState<string>(props.filePath || '');
  
  const handleFilePathChange = (value: string) => {
    setFilePath(value);
    props.onChange({ ...props, filePath: value });
  };
  
  return (
    <>
      <InlineField 
        labelWidth={LeftColumnWidth * 2} 
        label="File Path (optional)"
        tooltip="Optionally specify a file or directory to find owners for (e.g., 'src/main.go', 'src/pkg/services/'). Leave empty to show all CODEOWNERS entries for a repository."
      >
        <Input
          aria-label="File Path"
          width={RightColumnWidth}
          value={filePath}
          placeholder="e.g., src/main.go or src/pkg/services/"
          onChange={(el) => setFilePath(el.currentTarget.value)}
          onBlur={(el) => handleFilePathChange(el.currentTarget.value)}
        />
      </InlineField>
    </>
  );
};

export default QueryEditorCodeowners; 
