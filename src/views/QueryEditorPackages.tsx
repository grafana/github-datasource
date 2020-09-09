import React from 'react';

import { Input, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

import { QueryInlineField } from '../components/Forms';
import { PackagesOptions, PackageType } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends PackagesOptions {
  onChange: (value: PackagesOptions) => void;
}

const DefaultPackageType = PackageType.NPM;

const packageTypeOptions: Array<SelectableValue<string>> = Object.keys(PackageType).map(v => {
  return {
    label: v.replace('/_/gi', ' '),
    value: v,
  };
});

export default (props: Props) => {
  return (
    <>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Package type"
      >
        <Select
          options={packageTypeOptions}
          value={props.packageType || DefaultPackageType}
          width={RightColumnWidth}
          onChange={opt =>
            props.onChange({
              ...props,
              packageType: opt.value as PackageType,
            })
          }
        />
      </QueryInlineField>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Names"
        tooltip="Search for packages by their names"
      >
        <Input
          css=""
          value={props.query}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={el =>
            props.onChange({
              ...props,
              query: el.currentTarget.value,
            })
          }
        />
      </QueryInlineField>
    </>
  );
};
