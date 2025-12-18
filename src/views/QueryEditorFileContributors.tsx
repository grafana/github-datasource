import React from 'react';
import { Field, Input } from '@grafana/ui';

import { FileContributorsOptions } from '../types/query';

export interface QueryEditorFileContributorsProps {
  options: FileContributorsOptions;
  onOptionsChange: (options: FileContributorsOptions) => void;
}

export const QueryEditorFileContributors: React.FC<QueryEditorFileContributorsProps> = ({
  options,
  onOptionsChange,
}) => {
  const onFilePathChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      filePath: e.target.value,
    });
  };

  const onLimitChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const limit = parseInt(e.target.value, 10);
    onOptionsChange({
      ...options,
      limit: isNaN(limit) ? 10 : limit,
    });
  };

  return (
    <>
      <Field
        label="File Path"
        description="Enter the path to the file (e.g., src/components/Button.tsx)"
        required
      >
        <Input
          placeholder="src/components/Button.tsx"
          value={options.filePath || ''}
          onChange={onFilePathChange}
        />
      </Field>
      <Field
        label="Limit"
        description="Maximum number of contributors to return (default: 10)"
      >
        <Input
          type="number"
          placeholder="10"
          value={options.limit || 10}
          onChange={onLimitChange}
          min={1}
          max={100}
        />
      </Field>
    </>
  );
}; 