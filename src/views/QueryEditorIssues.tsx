import React, { useState } from 'react';
import { Combobox, ComboboxOption, Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import { components } from 'components/selectors';
import { IssueTimeField } from '../constants';
import type { IssuesOptions } from '../types/query';

interface Props extends IssuesOptions {
  onChange: (value: IssuesOptions) => void;
}

const timeFieldOptions: Array<ComboboxOption<IssueTimeField>> = Object.keys(IssueTimeField)
  .filter((_, i) => IssueTimeField[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${IssueTimeField[i]}`,
      value: i as IssueTimeField,
    };
  });

const defaultTimeField = 0 as IssueTimeField;

export const QueryEditorIssues = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <EditorRow>
      <EditorField
        label="Query"
        tooltip={() => (
          <>
            For more information, visit&nbsp;
            <a
              href="https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests"
              target="_blank"
              rel="noreferrer"
            >
              https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests
            </a>
          </>
        )}
      >
        <Input
          value={query}
          width={RightColumnWidth * 2 + 2}
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              query: el.currentTarget.value,
            })
          }
        />
      </EditorField>
      <EditorField label="Time Field" tooltip="The time field to filter on the time range">
        <Combobox
          data-testid={components.QueryEditor.Issues.timeFieldInput}
          width={RightColumnWidth}
          options={timeFieldOptions}
          value={props.timeField || defaultTimeField}
          onChange={(opt) =>
            props.onChange({
              ...props,
              timeField: opt.value,
            })
          }
        />
      </EditorField>
    </EditorRow>
  );
};
