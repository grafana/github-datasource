import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RepositoriesOptions } from '../types/query';
import { LeftColumnWidth, RightColumnWidth } from './QueryEditor';

interface Props extends RepositoriesOptions {
  onChange: (value: RepositoriesOptions) => void;
}

export const QueryEditorRepositories = ({
  propertyName = '',
  propertyValue = '',
  onChange,
}: Props) => {
  const [name, setName] = useState(propertyName);
  const [value, setValue] = useState(propertyValue);

  const update = () => onChange({ propertyName: name, propertyValue: value });

  return (
    <EditorRow>
      <EditorField
        label="Property name"
        tooltip="Optional custom property used to filter enriched organization repositories."
        width={LeftColumnWidth}
      >
        <Input
          aria-label="Repository custom property name"
          value={name}
          onChange={(event) => setName(event.currentTarget.value)}
          onBlur={() => update()}
          width={RightColumnWidth}
        />
      </EditorField>
      <EditorField
        label="Property value"
        tooltip="Exact property value to include. Use * or leave empty to include all repositories."
        width={LeftColumnWidth}
      >
        <Input
          aria-label="Repository custom property value"
          value={value}
          onChange={(event) => setValue(event.currentTarget.value)}
          onBlur={() => update()}
          width={RightColumnWidth}
        />
      </EditorField>
    </EditorRow>
  );
};
