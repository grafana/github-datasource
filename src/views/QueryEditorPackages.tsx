import React, { useEffect, useMemo, useState } from 'react';
import { Input, Select, InlineField } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { PackageType } from '../constants';
import type { PackagesOptions } from '../types/query';

interface Props extends PackagesOptions {
  onChange: (value: PackagesOptions) => void;
}

export const DefaultPackageType = PackageType.DOCKER;

const QueryEditorPackages = (props: Props) => {
  const [names, setNames] = useState<string>(props.names || '');

  // Set default package type if not set
  useEffect(() => {
    if (!props.packageType) {
      props.onChange({
        ...props,
        packageType: DefaultPackageType,
      });
    }
  }, [props]);

  const packageTypeOptions = useMemo(() => {
    // @TODO: These package types are not supported through GraphQL endpoint that we are using in our queries.
    // We should remove them from the list of options, but for now we will just ignore them
    // and not show them in the dropdown, if they have not been selected before.
    // Not sure if they ever been supported.
    const notSupportedTroughGraphQL: string[] = [PackageType.NPM, PackageType.RUBYGEMS, PackageType.NUGET];
    const packageTypeOptions: SelectableValue[] = Object.values(PackageType)
      .filter((packageType) => {
        // Filter out package types that are not supported through GraphQL
        return !notSupportedTroughGraphQL.includes(packageType);
      })
      .map((v) => {
        return {
          label: v.replace('/_/gi', ' '),
          value: v,
        };
      });

    // If user has selected a package type that is not in the list of options, add it to the list
    if (props.packageType) {
      const selectedPackageType = packageTypeOptions.find((opt: SelectableValue) => opt.value === props.packageType);
      if (!selectedPackageType) {
        packageTypeOptions.push({
          label: props.packageType.replace('/_/gi', ' '),
          value: props.packageType,
        });
      }
    }
    return packageTypeOptions;
    // We want to run this only once when component is mounted and not every time packageType is changed
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <InlineField labelWidth={LeftColumnWidth * 2} label="Package type">
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
      </InlineField>
      <InlineField
        labelWidth={LeftColumnWidth * 2}
        label="Names"
        tooltip="Search for packages using a comma delimited list of names"
      >
        <Input
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
      </InlineField>
    </>
  );
};

export default QueryEditorPackages;
