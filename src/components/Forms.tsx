import React, { InputHTMLAttributes, FunctionComponent } from 'react';
import { InlineFieldRow, InlineFormLabel } from '@grafana/ui';
import { css } from '@emotion/css';

export interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  tooltip?: string;
  labelWidth?: number;
  children?: React.ReactNode;
}

export const QueryField: FunctionComponent<Partial<Props>> = ({ label, labelWidth = 8, tooltip, children }) => (
  <>
    <InlineFormLabel width={labelWidth} className="query-keyword" tooltip={tooltip}>
      {label}
    </InlineFormLabel>
    {children}
  </>
);

const terminatorCss = css`
  margin-left: 4px;
`;

export const QueryRowTerminator = () => {
  return (
    <div className={`${terminatorCss} gf-form gf-form--grow`}>
      <div className="gf-form-label gf-form-label--grow" />
    </div>
  );
};

export const QueryInlineField = ({ ...props }) => {
  return (
    <QueryEditorRow>
      <QueryField {...props} />
    </QueryEditorRow>
  );
};

export const QueryEditorRow = (props: any) => {
  return (
    <InlineFieldRow>
      {props.children}
      <QueryRowTerminator />
    </InlineFieldRow>
  );
};
