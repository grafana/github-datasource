import React, { useState } from 'react';
import { Input, Combobox, ComboboxOption, MultiCombobox } from '@grafana/ui';
import { EditorField, EditorRow, QueryOptionGroup } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import { PullRequestFields, PullRequestTimeField, hasAllPullRequestFields } from '../constants';
import type { PullRequestField, Pull_RequestsOptions } from '../types/query';

interface Props extends Pull_RequestsOptions {
  onChange: (value: Pull_RequestsOptions) => void;
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

const fieldLabels: Record<PullRequestField, Omit<ComboboxOption<PullRequestField>, 'value'>> = {
  diffstat: { label: 'Additions / deletions', description: 'Slowest' },
  mergeable: { label: 'Mergeable' },
  author: { label: 'Author' },
  merged_by: { label: 'Merged by' },
  labels: { label: 'Labels', description: 'Up to 100 per pull request' },
};
const fieldOptions: Array<ComboboxOption<PullRequestField>> = PullRequestFields.map((value) => ({
  value,
  ...fieldLabels[value],
}));

const fieldsSummary = (fields?: PullRequestField[]) => {
  if (fields === undefined || hasAllPullRequestFields(fields)) {
    return 'Additional fields: All';
  }
  if (fields.length === 0) {
    return 'Additional fields: None';
  }
  return `Additional fields: ${fieldOptions
    .filter((o) => fields.includes(o.value))
    .map((o) => o.label)
    .join(', ')}`;
};

export const QueryEditorPullRequests = (props: Props) => {
  const [query, setQuery] = useState<string>(props.query || '');
  return (
    <>
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
      <EditorRow>
        <QueryOptionGroup title="Options" collapsedInfo={[fieldsSummary(props.fields)]}>
          <EditorField
            label="Additional fields"
            tooltip="These fields are slow or large for GitHub to return, so clearing them makes queries over many pull requests several times faster. All other columns are always included; use the Organize fields transformation to hide them."
          >
            <MultiCombobox
              aria-label="Additional fields"
              width="auto"
              minWidth={RightColumnWidth}
              options={fieldOptions}
              value={props.fields ?? PullRequestFields}
              enableAllOption
              onChange={(selected) => {
                const fields = selected.map((o) => o.value);
                // Leave fields unset when everything is selected, so new field groups are fetched by default
                props.onChange({ ...props, fields: hasAllPullRequestFields(fields) ? undefined : fields });
              }}
            />
          </EditorField>
        </QueryOptionGroup>
      </EditorRow>
    </>
  );
};
