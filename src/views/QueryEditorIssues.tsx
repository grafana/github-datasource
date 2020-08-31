import React, { ReactNode, useEffect, useState, useCallback } from 'react';

import { Select, Input, Tab, TabsBar, TabContent } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

import { ToggleField } from '../components/ToggleField';
import { QueryInlineField, QueryEditorRow } from '../components/Forms';
import LabelSelector from '../components/LabelSelector';
import { Label, IssuesOptions, IssueTimeField, IssueFilters } from '../types';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import { DataSource } from '../DataSource';

interface Props extends IssuesOptions {
  owner: string;
  repository: string;
  datasource: DataSource;
  onChange: (value: IssuesOptions) => void;
}

const timeFieldOptions: Array<SelectableValue<IssueTimeField>> = Object.keys(IssueTimeField)
  .filter((_, i) => IssueTimeField[i] !== undefined)
  .map((_, i) => {
    return {
      label: `${IssueTimeField[i]}`,
      value: i as IssueTimeField,
    };
  });

const TimeFieldSelect = (props: Props) => {
  return (
    <QueryInlineField
      labelWidth={LeftColumnWidth}
      label="Time Field"
      tooltip="The time field to filter on th time range"
    >
      <Select
        width={RightColumnWidth}
        options={timeFieldOptions}
        value={props.timeField || IssueTimeField.CreatedAt}
        onChange={opt =>
          props.onChange({
            ...props,
            timeField: opt.value,
          })
        }
      />
    </QueryInlineField>
  );
};

interface ToggleInputFieldProps {
  onChange: (value: IssuesOptions) => void;
  options: IssuesOptions;
  label: string;
  tooltip?: string;
  filter: keyof IssueFilters;
}

export const ToggleInputField = (props: ToggleInputFieldProps) => {
  const [value, setValue] = useState<string>();

  const { options, filter: key } = props;
  const { filters } = options;

  const disable = useCallback(() => {
    props.onChange({
      ...options,
      filters: {
        ...filters,
        [key]: undefined,
      },
    });
  }, [props.onChange, key]);

  const enable = useCallback(() => {
    props.onChange({
      ...options,
      filters: {
        ...filters,
        [key]: value || '',
      },
    });
  }, [props.onChange, key]);

  const inputDisabled = filters === undefined || filters[key] === undefined;
  return (
    <ToggleField
      checked={!inputDisabled}
      onToggle={() => (inputDisabled ? enable() : disable())}
      labelWidth={LeftColumnWidth}
      label={props.label}
      tooltip={props.tooltip}
    >
      <Input
        css=""
        disabled={inputDisabled}
        value={filters && filters[key]}
        width={RightColumnWidth}
        onChange={el => {
          props.onChange({
            ...options,
            filters: {
              ...filters,
              [key]: el.currentTarget.value,
            },
          });
          setValue(el.currentTarget.value);
        }}
      />
    </ToggleField>
  );
};

interface FieldProps extends Props {
  labels?: Label[];
  labelOptions?: Label[];
  onLabelsChange: (labels: Label[]) => void;
}

const Fields = (props: FieldProps) => {
  const { labels, labelOptions, onLabelsChange } = props;
  return (
    <>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Labels"
        tooltip="Selecting two or more labels will search for issues that have ALL labels, not either"
      >
        <LabelSelector
          width={RightColumnWidth * 2 + LeftColumnWidth}
          options={labelOptions || []}
          value={labels}
          onChange={onLabelsChange}
        />
      </QueryInlineField>
      <QueryEditorRow>
        <ToggleInputField label="Assignee" filter="assignee" options={props} onChange={props.onChange} />
        <ToggleInputField label="Author" filter="createdBy" options={props} onChange={props.onChange} />
      </QueryEditorRow>
      <QueryEditorRow>
        <ToggleInputField label="Mentioned" filter="mentioned" options={props} onChange={props.onChange} />
        <ToggleInputField label="Milestone" filter="milestone" options={props} onChange={props.onChange} />
      </QueryEditorRow>
    </>
  );
};

const Query = (props: Props) => {
  return (
    <>
      <QueryInlineField
        labelWidth={LeftColumnWidth}
        label="Query"
        tooltip="For more information, visit https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests"
      >
        <Input
          css=""
          value={props.query}
          width={RightColumnWidth * 2 + LeftColumnWidth}
          onChange={el =>
            props.onChange({
              ...props,
              query: el.currentTarget.value,
            })
          }
        />
      </QueryInlineField>
      <TimeFieldSelect {...props} />
    </>
  );
};

interface TabListItem {
  label: string;
  key: string;
  component: (props: FieldProps) => ReactNode;
}

const tabs: TabListItem[] = [
  { label: 'Fields', key: 'fields', component: props => <Fields {...props} /> },
  { label: 'Query', key: 'query', component: props => <Query {...props} /> },
];

export default (props: Props) => {
  const [activeTab, setActiveTab] = useState<TabListItem>(tabs[0]);

  const [labels, setLabels] = useState<Label[]>();
  const [labelOptions, setLabelOptions] = useState<Label[]>();

  useEffect(() => {
    const res = props.datasource.getLabels(props.repository || '', props.owner || '');
    res.then(labels => setLabelOptions(labels));
    res.catch(err => console.error(err));
  }, [props.repository, props.owner]);

  const onLabelsChange = useCallback(
    (labels: Label[]) => {
      setLabels(labels);
      // Undefined = no labels, whereas an empty Array of labels will have 0 results
      props.onChange({
        ...props,
        filters: {
          ...props.filters,
          labels: labels.length > 0 ? labels?.map(label => label.name) : undefined,
        },
      });
    },
    [props.onChange]
  );

  return (
    <>
      <TabsBar>
        {tabs.map((tab, index) => {
          return (
            <Tab
              css=""
              key={index}
              label={tab.label}
              active={tab === activeTab}
              tabIndex={index}
              onChangeTab={() => setActiveTab(tab)}
            />
          );
        })}
      </TabsBar>
      <TabContent>
        {activeTab.component({
          ...props,
          onLabelsChange,
          labels,
          labelOptions,
        })}
      </TabContent>
    </>
  );
};
