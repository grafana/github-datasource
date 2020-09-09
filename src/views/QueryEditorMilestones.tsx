import React from 'react';

import { Input } from '@grafana/ui';

import { QueryInlineField } from '../components/Forms';
import { MilestonesOptions } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends MilestonesOptions {
  onChange: (value: MilestonesOptions) => void;
}

export default (props: Props) => {
  return (
    <>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Query"
        tooltip="Query milestones by title"
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
