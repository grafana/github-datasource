import React, { useState } from 'react';

import { Input, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

import { QueryInlineField } from '../components/Forms';
import { PackagesOptions, PackageType } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends PackagesOptions {
  onChange: (value: PackagesOptions) => void;
}

const DefaultPackageType = PackageType.NPM;

const packageTypeOptions: Array<SelectableValue<string>> = Object.keys(PackageType).map((v) => {
  return {
    label: v.replace('/_/gi', ' '),
    value: v,
  };
});

const QueryEditorPackages = (props: Props) => {
  const [names, setNames] = useState<string>(props.names || '');

  return (
    <>
      <QueryInlineField labelWidth={LeftColumnWidth} label="Package type">
        <Select
          options={packageTypeOptions}
          value={props.packageType || DefaultPackageType}
          width={RightColumnWidth}
          onChange={(opt) =>
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
        tooltip="Search for packages using a comma delimited list of names"
      >
        <Input
          css=""
          value={names}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={(el) => setNames(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props.query,
              names: el.currentTarget.value,
            })
          }
        />
      </QueryInlineField>
    </>
  );
};

export default QueryEditorPackages;
