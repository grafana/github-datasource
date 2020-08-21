import React from 'react';
import { Input } from '@grafana/ui';
import QueryEditorRepository from './QueryEditorRepository';

import { QueryInlineField } from '../components/Forms';
import { CommitsOptions } from '../query';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';

interface Props extends CommitsOptions {
  onChange: (value: CommitsOptions) => void;
};

export default (props: Props) => {
  return (
    <>
      <QueryEditorRepository {...props} />
      <QueryInlineField labelWidth={LeftColumnWidth} label='Ref (Branch / Tag)'>
        <Input
          css=''
          width={RightColumnWidth}
          value={props.gitRef}
          onChange={el => props.onChange({...props, gitRef: el.currentTarget.value})}
        />
      </QueryInlineField>
    </>
  );
}
