import React, { InputHTMLAttributes, FunctionComponent } from 'react';
import { InlineFormLabel } from '@grafana/ui';

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

export const QueryInlineField = ({ ...props }) => {
  return (
    <div className={'gf-form-inline'}>
      <QueryField {...props} />
      <div className="gf-form gf-form--grow">
        <div className="gf-form-label gf-form-label--grow" />
      </div>
    </div>
  );
};

export const QueryEditorRow = (props: any) => {
  return <div className="gf-form-inline">{props.children}</div>;
};
