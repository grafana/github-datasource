import React from 'react';
import { InlineFormLabel, Switch } from '@grafana/ui';
import { css } from 'emotion';

const switchStyle = css`
line-height: inherit;
`;

const flexFixStyle = css`
margin: 0px;
`;

const containerStyle = css`
margin-right: 4px;
`;

const switchWidth = 3;

type Props = {
  checked: boolean;
  onToggle: () => void;
  label: string;
  labelWidth: number;
  tooltip?: string;
  children: any;
}

export const ToggleField = (props: Props) => {
  return (
    <div className='gf-form-inline'>
      <div className={`gf-form-inline ${containerStyle} width-${props.labelWidth}`}>
        <div>
          <InlineFormLabel width={props.labelWidth - switchWidth} className={`query-keyword ${flexFixStyle}`} tooltip={props.tooltip}>
            {props.label}
          </InlineFormLabel>
        </div>
        <div>
          <div className={`gf-form-label ${flexFixStyle} ${switchStyle} width-${switchWidth}`}>
            <Switch css='' checked={props.checked} onChange={props.onToggle} />
          </div>
        </div>
      </div>
      {props.children}
    </div>
  );
};
