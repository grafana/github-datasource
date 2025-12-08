import React, { useState } from 'react';
import { Input, InlineField } from '@grafana/ui';
import { RightColumnWidth, LeftColumnWidth } from './QueryEditor';
import type { TeamsOptions } from '../types/query';

interface Props extends TeamsOptions {
  organization?: string;
  onChange: (value: TeamsOptions & { organization?: string }) => void;
}

const QueryEditorTeams = (props: Props) => {
  const [organization, setOrganization] = useState<string>(props.organization || '');
  const [query, setQuery] = useState<string>(props.query || '');
  
  const handleOrganizationChange = (value: string) => {
    setOrganization(value);
    props.onChange({ ...props, organization: value });
  };
  
  const handleQueryChange = (value: string) => {
    setQuery(value);
    props.onChange({ ...props, query: value });
  };
  
  return (
    <>
      <InlineField 
        labelWidth={LeftColumnWidth * 2} 
        label="Organization"
        tooltip="GitHub organization to query teams for (e.g., 'grafana')"
      >
        <Input
          aria-label="Organization"
          width={RightColumnWidth}
          value={organization}
          placeholder="e.g., grafana"
          onChange={(el) => setOrganization(el.currentTarget.value)}
          onBlur={(el) => handleOrganizationChange(el.currentTarget.value)}
        />
      </InlineField>
      
      <InlineField 
        labelWidth={LeftColumnWidth * 2} 
        label="Query (optional)"
        tooltip="Filter teams by name/description, or fetch a specific team by slug. Examples: 'backend' (filters), 'platform-team' (specific team slug). Leave empty to show all teams."
      >
        <Input
          aria-label="Team Query"
          width={RightColumnWidth}
          value={query}
          placeholder="e.g., backend or platform-team"
          onChange={(el) => setQuery(el.currentTarget.value)}
          onBlur={(el) => handleQueryChange(el.currentTarget.value)}
        />
      </InlineField>
    </>
  );
};

export default QueryEditorTeams; 