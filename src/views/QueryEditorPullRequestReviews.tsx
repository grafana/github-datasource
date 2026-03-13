import React, { useState } from 'react';
import { Input, Combobox, ComboboxOption } from '@grafana/ui';
import { RightColumnWidth } from './QueryEditor';
import { PullRequestTimeField } from '../constants';
import type { PullRequestReviewsOptions } from '../types/query';
import { EditorField, EditorRow } from '@grafana/plugin-ui';

interface Props extends PullRequestReviewsOptions {
  onChange: (value: PullRequestReviewsOptions) => void;
}

const timeFieldOptions: Array<ComboboxOption<PullRequestTimeField>> = Object.keys(PullRequestTimeField)
  .filter((_, i) => PullRequestTimeField[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${PullRequestTimeField[i]}`,
      value: i as PullRequestTimeField,
    };
  });

const defaultTimeField = timeFieldOptions[0].value;

export const QueryEditorPullRequestReviews = (props: Props) => {
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
      <EditorField
        label="Time Field"
        tooltip="The time field to filter on the time range. WARNING: If selecting 'None', be mindful of the amount of data being queried. On larger repositories, querying all pull requests could easily cause rate limiting"
      >
        <Combobox
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
