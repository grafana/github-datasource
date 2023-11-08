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

export const QueryRowTerminator = () => {
  const styles = getStyles();

  return (
    <InlineFormLabel className={styles.rowTerminator}>
      <></>
    </InlineFormLabel>
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
  const styles = getStyles();

  return (
    <InlineFieldRow className={styles.rowSpacing}>
      {props.children}
      <QueryRowTerminator />
    </InlineFieldRow>
  );
};

const getStyles = () => {
  return {
    rowSpacing: css({
      marginBottom: '4px',
    }),
    rowTerminator: css({
      flexGrow: 1,
      marginLeft: '4px',
    }),
  };
};
