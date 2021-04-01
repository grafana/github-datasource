import React, { useMemo, useState, useCallback } from 'react';
import coreModule from 'grafana/app/core/core_module';

import { AnnotationQueryRequest } from '@grafana/data';

import QueryEditor from './QueryEditor';
import { DataSource } from '../DataSource';
import { GitHubAnnotationQuery, QueryType } from '../types';
import { QueryInlineField } from '../components/Forms';
import FieldSelect from '../components/FieldSelect';
import { isValid } from '../validation';
import { selectors } from 'components/selectors';

interface Props {
  datasource: DataSource;
  annotation: AnnotationQueryRequest<GitHubAnnotationQuery>;
  change: (query: AnnotationQueryRequest<GitHubAnnotationQuery>) => void;
}

const AnnotationQueryEditor = (props: Props) => {
  const [choices, setChoices] = useState<string[]>();
  const { annotation } = props;

  useMemo(async () => {
    if (isValid(props.annotation.annotation)) {
      setChoices(await props.datasource.getChoices(props.annotation.annotation));
    }
  }, [props.annotation.annotation, props.datasource]);

  const onChange = useCallback(
    (query: GitHubAnnotationQuery) => {
      props.change({
        ...props.annotation,
        annotation: {
          ...annotation.annotation,
          ...query,
          datasource: props.datasource.name,
        },
      });
    },
    [props, annotation.annotation]
  );

  return (
    <div aria-label={selectors.components.AnnotationEditor.container}>
      <QueryEditor
        query={annotation.annotation}
        datasource={props.datasource}
        onChange={(query) =>
          onChange({
            ...query,
            field: annotation.annotation.field,
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
      {isValid(props.annotation.annotation) && (
        <>
          <QueryInlineField
            width={10}
            labelWidth={10}
            label="Display Field"
            tooltip="This field determines the text / value displayed on the annotation"
          >
            <FieldSelect
              onChange={(value) =>
                onChange({
                  ...annotation.annotation,
                  field: value,
                })
              }
              options={choices || []}
              width={64}
              value={annotation.annotation.field}
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
                onChange({
                  ...annotation.annotation,
                  timeField: value,
                })
              }
              options={choices || []}
              width={64}
              value={annotation.annotation.timeField}
              loading={!choices}
            />
          </QueryInlineField>
        </>
      )}
    </div>
  );
};

coreModule.directive('githubAnnotationEditor', [
  'reactDirective',
  (reactDirective: any) => {
    return reactDirective(AnnotationQueryEditor, ['annotation', 'datasource', 'change']);
  },
]);
