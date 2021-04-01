import React from 'react';
import { Label } from '../types';
import { MultiSelect } from '@grafana/ui';

interface Props {
  options: Label[];
  value?: Label[];
  onChange: (labels: Label[]) => void;
  width?: number;
}

const LabelSelector = (props: Props) => {
  return (
    <MultiSelect
      width={props.width}
      options={props.options.map((value) => {
        return {
          label: value.name,
          value: value,
        };
      })}
      value={props.value}
      onChange={(values) => props.onChange(values?.map((item) => item.value!))}
    />
  );
};

export default LabelSelector;
