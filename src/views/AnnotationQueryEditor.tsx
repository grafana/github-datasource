import React, { useMemo, useCallback, useState } from 'react';
import { DataSourceJsonData, QueryEditorProps } from '@grafana/data';

import QueryEditor from './QueryEditor';
import { GithubDataSource } from '../DataSource';
import { GitHubAnnotationQuery, GitHubQuery, QueryType } from '../types';
import { QueryInlineField } from '../components/Forms';
import FieldSelect from '../components/FieldSelect';
import { isValid } from '../validation';
import { selectors } from 'components/selectors';

export const AnnotationQueryEditor = (
  props: QueryEditorProps<GithubDataSource, GitHubAnnotationQuery, DataSourceJsonData>
) => {
  const [choices, setChoices] = useState<string[]>();
  const { query, onChange } = props;

  useMemo(async () => {
    if (isValid(query as unknown as GitHubQuery)) {
      setChoices(await props.datasource.getChoices(query as unknown as GitHubQuery));
    }
  }, [query, props.datasource]);

  const handleDnChange = useCallback(
    (annotationQuery: GitHubAnnotationQuery) => {
      onChange({
        ...query,
        ...annotationQuery,
      });
    },
    [props, query]
  );

  return (
    <div aria-label={selectors.components.AnnotationEditor.container}>
      <QueryEditor
        query={query as unknown as GitHubQuery}
        datasource={props.datasource}
        onChange={(query) =>
          handleDnChange({
            ...query,
            field: query.field,
          })
        }
        onRunQuery={() => {}}
        queryTypes={[
          QueryType.Commits,
          QueryType.Releases,
          QueryType.Pull_Requests,
          QueryType.Issues,
          QueryType.Milestones,
          QueryType.Tags,
        ]}
      />

      {/* Only display the field selection items when the user has created an actual query */}
      {isValid(query as unknown as GitHubQuery) && (
        <>
          <QueryInlineField
            width={10}
            labelWidth={10}
            label="Display Field"
            tooltip="This field determines the text / value displayed on the annotation"
          >
            <FieldSelect
              onChange={(value) =>
                handleDnChange({
                  ...query,
                  field: value,
                } as unknown as GitHubQuery)
              }
              options={choices || []}
              width={64}
              value={query.field}
              loading={!choices}
            />
          </QueryInlineField>
          <QueryInlineField
            width={10}
            labelWidth={10}
            label="Time Field"
            tooltip="This field is used to determine where the annotation will display on a graph"
          >
            <FieldSelect
              onChange={(value) =>
                handleDnChange({
                  ...query,
                  timeField: value,
                } as unknown as GitHubQuery)
              }
              options={choices || []}
              width={64}
              value={query.timeField}
              loading={!choices}
            />
          </QueryInlineField>
        </>
      )}
    </div>
  );
};
