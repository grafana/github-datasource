import React from 'react';
import { Input, InlineField, DateTimePicker } from '@grafana/ui';
import { dateTime } from '@grafana/data';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { CopilotMetricsTeamOptions } from '../types/query';

interface Props extends CopilotMetricsTeamOptions {
  onChange: (val: CopilotMetricsTeamOptions) => void;
}

const QueryEditorCopilotMetricsTeam = (props: Props) => {
  return (
    <>
      <InlineField labelWidth={LeftColumnWidth * 2} label="Organization" tooltip="GitHub organization name">
        <Input
          value={props.organization || ''}
          width={RightColumnWidth}
          onChange={(e) =>
            props.onChange({
              ...props,
              organization: e.currentTarget.value,
            })
          }
          placeholder="Enter organization name"
        />
      </InlineField>
      
      <InlineField labelWidth={LeftColumnWidth * 2} label="Team Slug" tooltip="GitHub team slug name">
        <Input
          value={props.teamSlug || ''}
          width={RightColumnWidth}
          onChange={(e) =>
            props.onChange({
              ...props,
              teamSlug: e.currentTarget.value,
            })
          }
          placeholder="Enter team slug"
        />
      </InlineField>
      
      <InlineField 
        labelWidth={LeftColumnWidth * 2} 
        label="Since" 
        tooltip="Show usage metrics since this date (max 28 days ago)"
      >
        <DateTimePicker
          date={props.since ? dateTime(props.since) : undefined}
          onChange={(date) =>
            props.onChange({
              ...props,
              since: date?.toISOString(),
            })
          }
        />
      </InlineField>
      
      <InlineField 
        labelWidth={LeftColumnWidth * 2} 
        label="Until" 
        tooltip="Show usage metrics until this date"
      >
        <DateTimePicker
          date={props.until ? dateTime(props.until) : undefined}
          onChange={(date) =>
            props.onChange({
              ...props,
              until: date?.toISOString(),
            })
          }
        />
      </InlineField>
      
      <InlineField 
        labelWidth={LeftColumnWidth * 2} 
        label="Page" 
        tooltip="Page number for pagination"
      >
        <Input
          type="number"
          value={props.page || ''}
          width={RightColumnWidth / 2}
          onChange={(e) =>
            props.onChange({
              ...props,
              page: parseInt(e.currentTarget.value, 10) || undefined,
            })
          }
          placeholder="1"
          min={1}
        />
      </InlineField>
      
      <InlineField 
        labelWidth={LeftColumnWidth * 2} 
        label="Per Page" 
        tooltip="Number of days of metrics to display per page (max 28)"
      >
        <Input
          type="number"
          value={props.perPage || ''}
          width={RightColumnWidth / 2}
          onChange={(e) =>
            props.onChange({
              ...props,
              perPage: parseInt(e.currentTarget.value, 10) || undefined,
            })
          }
          placeholder="28"
          min={1}
          max={28}
        />
      </InlineField>
    </>
  );
};

export default QueryEditorCopilotMetricsTeam;
