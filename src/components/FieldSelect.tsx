import React from 'react';
import { Select, Spinner } from '@grafana/ui';
import { css } from 'emotion';

interface Props {
  onChange: (value: string) => void;
  value?: string;
  options: string[];
  width: number;
  loading: boolean;
}

const containerCss = css`
  align-items: center;
`;

const spannerCss = css`
  margin: 0px 3px;
  padding: 0px 3px;
`;

const FieldSelect = (props: Props) => {
  const { onChange, options, value, width, loading } = props;
  return (
    <div className={`${containerCss} gf-form-inline`}>
      <Select
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
    </div>
  );
};

export default FieldSelect;
