import React, { useMemo, useState } from 'react';

import { Select } from '@grafana/ui';

import QueryEditor from './QueryEditor';
import { DataSource } from '../DataSource';
import { GitHubVariableQuery, DefaultQueryType, QueryType } from '../types';
import { QueryInlineField } from '../components/Forms';
import { isValid } from '../validation';

interface Props {
  query: GitHubVariableQuery;
  onChange: (query: GitHubVariableQuery, definition: string) => void;
  datasource: DataSource;
}

export default (props: Props) => {
  const definition = `${props.datasource.name} - ${props.query.queryType || DefaultQueryType}`;
  const [choices, setChoices] = useState<string[]>();

  useMemo(async () => {
    if (isValid(props.query)) {
      setChoices(await props.datasource.getChoices(props.query));
    }
  }, [props.query]);

  return (
    <>
      <QueryEditor
        query={props.query}
        datasource={props.datasource}
        onChange={query =>
          props.onChange(
            {
              ...query,
              field: props.query.field,
            },
            definition
          )
        }
        onRunQuery={() => {}}
        queryTypes={[QueryType.Contributors, QueryType.Tags, QueryType.Releases, QueryType.Labels]}
      />
      <QueryInlineField width={10} labelWidth={10} label="field">
        <Select
          width={64}
          options={choices?.map(v => {
            return { label: v, value: v };
          })}
          value={props.query.field}
          onChange={opt =>
            props.onChange(
              {
                ...props.query,
                field: opt.value,
              },
              definition
            )
          }
        />
      </QueryInlineField>
    </>
  );
};
