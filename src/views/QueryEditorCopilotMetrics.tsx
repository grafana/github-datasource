import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { CopilotMetricsOptions } from '../types/query';
import { components } from '../components/selectors';

interface Props extends CopilotMetricsOptions {
  onChange: (val: CopilotMetricsOptions) => void;
}

const QueryEditorCopilotMetrics = (props: Props) => {
  const [team, setTeam] = useState<string>(props.teamSlug || '');

  return (
    <>
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
