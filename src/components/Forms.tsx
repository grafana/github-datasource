import React, { InputHTMLAttributes } from 'react';
import { InlineFieldRow } from '@grafana/ui';
import { css } from '@emotion/css';

export interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  tooltip?: string;
  labelWidth?: number;
  children?: React.ReactNode;
}

export const QueryEditorRow = (props: any) => {
  const styles = getStyles();
  return <InlineFieldRow className={styles.rowSpacing}>{props.children}</InlineFieldRow>;
};

const getStyles = () => {
  return {
    rowSpacing: css({
      marginBottom: '4px',
    }),
  };
};
