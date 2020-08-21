import React, { useEffect, useState, useCallback } from 'react';
import { Input } from '@grafana/ui';

import { QueryInlineField } from '../components/Forms';
import LabelSelector from '../components/LabelSelector';
import QueryEditorRepository from './QueryEditorRepository';
import { Label } from '../types';
import { IssuesOptions } from '../query';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { DataSource } from '../DataSource';

interface Props extends IssuesOptions {
  datasource: DataSource;
  onChange: (value: IssuesOptions) => void;
};

export default (props: Props) => {
  const { filters } = props;
  const [labels, setLabels] = useState<Label[]>();
  const [labelOptions, setLabelOptions] = useState<Label[]>();

  useEffect(
    () => {
      const res = props.datasource.getLabels(props.repository || '', props.owner || '');
      res.then(labels => setLabelOptions(labels));
      res.catch(err => console.error(err));
    },
    [props.repository, props.owner],
  );

  const onLabelsChange = useCallback((labels: Label[]) => {
    setLabels(labels);
    // Undefined = no labels, whereas an empty Array of labels will have 0 results
    props.onChange({
      ...props,
      filters: {
        ...props.filters,
        labels: labels.length > 0 ? labels?.map(label => label.name) : undefined,
      }
    });
  }, [props.onChange]);

  return (
    <>
      <QueryEditorRepository {...props} />
      <QueryInlineField labelWidth={LeftColumnWidth} label="Assignee">
        <Input css='' value={filters?.assignee} width={RightColumnWidth} onChange={el => props.onChange({
          ...props,
          filters: {
            ...filters,
            assignee: el.currentTarget.value,
          },
        })}/>
      </QueryInlineField>
      <QueryInlineField labelWidth={LeftColumnWidth} label="Labels" tooltip="Selecting two or more labels will search for issues that have ALL labels, not either">
        <LabelSelector width={RightColumnWidth * 2} options={labelOptions || []} value={labels} onChange={onLabelsChange} />
      </QueryInlineField>
    </>
  );
}
