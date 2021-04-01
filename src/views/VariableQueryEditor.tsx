import React, { useMemo, useState } from 'react';

import QueryEditor from './QueryEditor';
import { DataSource } from '../DataSource';
import { GitHubVariableQuery, DefaultQueryType, QueryType } from '../types';
import { QueryInlineField } from '../components/Forms';
import FieldSelect from '../components/FieldSelect';
import { isValid } from '../validation';

interface Props {
  query: GitHubVariableQuery;
  onChange: (query: GitHubVariableQuery, definition: string) => void;
  datasource: DataSource;
}

const VariableQueryEditor = (props: Props) => {
  const definition = `${props.datasource.name} - ${props.query.queryType || DefaultQueryType}`;
  const [choices, setChoices] = useState<string[]>();

  useMemo(async () => {
    if (isValid(props.query)) {
      setChoices(await props.datasource.getChoices(props.query));
    }
  }, [props.query, props.datasource]);

  return (
    <>
      <QueryEditor
        query={props.query}
        datasource={props.datasource}
        onChange={(query) =>
          props.onChange(
            {
              ...query,
              field: props.query.field,
            },
            definition
          )
        }
        onRunQuery={() => {}}
        queryTypes={[
          QueryType.Repositories,
          QueryType.Contributors,
          QueryType.Tags,
          QueryType.Releases,
          QueryType.Labels,
          QueryType.Milestones,
        ]}
      />
      <QueryInlineField
        width={10}
        labelWidth={10}
        label="Display Field"
        tooltip="This field determines the text / value used for the variable"
      >
        <FieldSelect
          onChange={(value) =>
            props.onChange(
              {
                ...props.query,
                field: value,
              },
              definition
            )
          }
          options={choices || []}
          width={64}
          value={props.query.field}
          loading={!choices}
        />
      </QueryInlineField>
    </>
  );
};

export default VariableQueryEditor;
