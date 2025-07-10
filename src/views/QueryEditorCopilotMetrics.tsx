import React, { useState } from 'react';
import { Input, InlineField, InlineLabel } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { CopilotMetricsOptions } from '../types/query';
import { components } from '../components/selectors';
import { QueryEditorRow } from '../components/Forms';

interface Props extends CopilotMetricsOptions {
  onChange: (val: CopilotMetricsOptions) => void;
}

const QueryEditorCopilotMetrics = (props: Props) => {
  const [team, setTeam] = useState(props.teamSlug || '');
  const [organization, setOrganization] = useState<string>(props.organization || '');

  return (
    <>
        <QueryEditorRow>
          <InlineLabel
            tooltip="The owner (organization or user) of the GitHub repository (example: 'grafana')"
            width={LeftColumnWidth * 2}
          >
            Owner
          </InlineLabel>
          <Input
            aria-label={components.QueryEditor.Owner.input}
            width={RightColumnWidth}
            value={organization}
            onChange={(el) => setOrganization(el.currentTarget.value)}
            onBlur={(el) =>
              props.onChange({
                ...props,
                organization: el.currentTarget.value,
              })
            }
          />
        </QueryEditorRow>
      <InlineField labelWidth={LeftColumnWidth * 2} label="Team Slug" tooltip="GitHub team slug name (optional - leave empty for organization-wide metrics)">
        <Input
          aria-label={components.QueryEditor.Owner.input}
          width={RightColumnWidth}
          value={team}
          onChange={(el) => setTeam(el.currentTarget.value)}
          onBlur={(el) =>
            props.onChange({
              ...props,
              teamSlug: el.currentTarget.value
            })
          }
          placeholder="Enter team slug (optional)"
        />
      </InlineField>
    </>
  );
};

export default QueryEditorCopilotMetrics;
