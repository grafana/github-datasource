import React, { useState } from 'react';
import { Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { CustomPropertiesOptions } from '../types/query';
import { LeftColumnWidth, RightColumnWidth } from './QueryEditor';

interface Props extends CustomPropertiesOptions {
  onChange: (value: CustomPropertiesOptions) => void;
}

export const QueryEditorCustomProperties = ({ propertyName = '', onChange }: Props) => {
  const [name, setName] = useState(propertyName);

  return (
    <EditorRow>
      <EditorField
        label="Property name"
        tooltip="Optional exact custom-property name. Leave empty to return the complete organization schema."
        width={LeftColumnWidth}
      >
        <Input
          aria-label="Custom property name"
          value={name}
          onChange={(event) => setName(event.currentTarget.value)}
          onBlur={() => onChange({ propertyName: name })}
          width={RightColumnWidth}
        />
      </EditorField>
    </EditorRow>
  );
};
