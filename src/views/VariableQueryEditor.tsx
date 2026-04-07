import React, { useEffect, useState } from 'react';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import QueryEditor from './QueryEditor';
import { GitHubDataSource } from '../DataSource';
import FieldSelect from '../components/FieldSelect';
import { isValid } from '../validation';
import { DefaultQueryType } from '../constants';
import type { GitHubVariableQuery } from '../types/query';

interface Props {
  query: GitHubVariableQuery;
  onChange: (query: GitHubVariableQuery, definition: string) => void;
  datasource: GitHubDataSource;
}

const VariableQueryEditor = (props: Props) => {
  const definition = `${props.datasource.name} - ${props.query.queryType || DefaultQueryType}`;
  const [choices, setChoices] = useState<string[]>();

  useEffect(() => {
    // Used to ignore stale responses when the query changes before a fetch completes (race condition)
    let ignore = false;

    async function fetchData() {
      if (isValid(props.query)) {
        const result = await props.datasource.getChoices(props.query);
        if (!ignore) {
          setChoices(result);
        }
      }
    }
    fetchData();

    return () => {
      ignore = true;
    };
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
      />
      <EditorRow>
        <EditorField label="Field Value" tooltip="This field determines the value used for the variable">
          <FieldSelect
            onChange={(value) =>
              props.onChange(
                {
                  ...props.query,
                  key: value,
                },
                definition
              )
            }
            options={choices || []}
            width={64}
            value={props.query.key}
            loading={!choices}
          />
        </EditorField>
        <EditorField label="Field Display" tooltip="This field determines the text used for the variable">
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
        </EditorField>
      </EditorRow>
    </>
  );
};

export default VariableQueryEditor;
