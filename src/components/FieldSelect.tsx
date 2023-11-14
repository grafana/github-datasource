import React from 'react';
import { InlineFieldRow, Select, Spinner } from '@grafana/ui';
import { css } from '@emotion/css';

interface Props {
  onChange: (value: string) => void;
  value?: string;
  options: string[];
  width: number;
  loading: boolean;
}

const spannerCss = css`
  margin: 0px 3px;
  padding: 0px 3px;
`;

const FieldSelect = (props: Props) => {
  const { onChange, options, value, width, loading } = props;
  return (
    <InlineFieldRow>
      <Select
        allowCustomValue
        value={value}
        onChange={(opt) => onChange(opt.value!)}
        width={width}
        disabled={loading}
        placeholder={loading ? 'Loading...' : 'Select...'}
        options={options.map((opt) => {
          return {
            label: opt,
            value: opt,
          };
        })}
      />
      {loading && (
        <div>
          <Spinner className={spannerCss} />
        </div>
      )}
    </InlineFieldRow>
  );
};

export default FieldSelect;
