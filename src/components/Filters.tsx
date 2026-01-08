import React, { useState } from 'react';
import { SegmentAsync, Segment, InlineFieldRow, Button } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { css } from '@emotion/css';

interface Props {
  value?: Filter[];
  onChange: (filters: Filter[]) => void;
  loadOptions: (filter?: string) => Promise<Array<SelectableValue<string>>>;
  ops: Array<SelectableValue<string>>;
}

export const Filters: React.FC<Props> = (props: Props) => {
  const { value, onChange } = props;
  const defaultFilters: Filter[] = [];
  const [filters, setFilters] = useState<Filter[]>(value || defaultFilters);
  const [loading, setLoading] = useState<number>();

  const add = () => {
    setFilters([...filters, { key: '', value: '', op: props.ops[0]?.value || '=' }]);
  };

  const onKeyChange = (index: number) => (selected: SelectableValue) => {
    const update = filters.map((f, i) => (i === index ? { ...f, key: selected.value, value: '' } : f));
    setFilters(update);
  };

  const onValueChange = (index: number) => (selected: SelectableValue) => {
    const update = filters.map((f, i) => (i === index ? { ...f, value: selected.value } : f));
    changeFilters(update);
  };

  const onOpChange = (index: number) => (selected: SelectableValue) => {
    const update = filters.map((f, i) => (i === index ? { ...f, op: selected.value } : f));
    changeFilters(update);
  };

  const onConjunctionChange = (index: number) => (selected: SelectableValue) => {
    const update = filters.map((f, i) => (i === index ? { ...f, conjunction: selected.value } : f));
    changeFilters(update);
  };

  const remove = (index: number) => () => {
    const update = filters.filter((f, i) => i !== index);
    changeFilters(update);
  };

  function changeFilters(filters: Filter[]) {
    setFilters(filters);
    onChange(filters);
  }

  const loadValues = (index: number) => async () => {
    setLoading(index);
    const key = filters[index].key;
    const opts = await props.loadOptions(key);
    setLoading(undefined);
    return opts;
  };

  const list = filters || [];

  const ops = ['and', 'or'];
  const opList = ops.map((op) => ({ label: op, value: op }));

  const styles = {
    loading: css`
      position: absolute;
      top: 0;
      padding: 8px 10px;
      font-weight: 500;
      font-size: 12px;
    `,
    wrapper: css`
      position: relative;
    `,
  };

  return (
    <InlineFieldRow>
      <>
        {list.map((filter: Filter, i: number) => (
          <>
            <SegmentAsync
              loadOptions={props.loadOptions}
              placeholder="Key..."
              allowCustomValue={true}
              value={filter.key}
              onChange={onKeyChange(i)}
            />
            <Segment placeholder="Operator..." value={filter.op} onChange={onOpChange(i)} options={props.ops} />
            <span className={styles.wrapper}>
              <SegmentAsync
                loadOptions={loadValues(i)}
                placeholder="Value..."
                allowCustomValue={true}
                value={filter.value}
                onChange={onValueChange(i)}
              />
              {loading === i && <div className={styles.loading}>Loading...</div>}
            </span>
            <Button variant="secondary" aria-label="Remove filter" onClick={remove(i)} icon="trash-alt" />
            {list.length > 1 && i !== list.length - 1 && (
              <Segment onChange={onConjunctionChange(i)} options={opList} value={filter.conjunction || 'and'}></Segment>
            )}
          </>
        ))}
        <Button variant="secondary" aria-label="Add filter" onClick={add} icon="plus" />
      </>
    </InlineFieldRow>
  );
};

Filters.displayName = 'Filters';

export interface Filter {
  key: string;
  value: string;
  op?: string;
  conjunction?: string;
}
