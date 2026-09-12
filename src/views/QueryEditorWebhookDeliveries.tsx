import React, { useState } from 'react';
import { Input, Combobox, ComboboxOption } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { RightColumnWidth } from './QueryEditor';
import { components } from '../components/selectors';
import { WebhookDeliveryStatus } from '../constants';
import type { WebhookDeliveriesOptions } from '../types/query';

interface Props extends WebhookDeliveriesOptions {
  onChange: (value: WebhookDeliveriesOptions) => void;
}

const statusOptions: ComboboxOption[] = [
  { label: 'All', value: WebhookDeliveryStatus.All },
  { label: 'Success', value: WebhookDeliveryStatus.Success },
  { label: 'Failure', value: WebhookDeliveryStatus.Failure },
];

export const QueryEditorWebhookDeliveries = (props: Props) => {
  const [hookId, setHookId] = useState<string>(props.hookId || '');
  const [event, setEvent] = useState<string>(props.event || '');

  return (
    <EditorRow>
      <EditorField
        label="Hook ID"
        tooltip="The numeric ID of the webhook. Leave the repository empty to query an organization webhook, or set it to query a webhook of that repository. Find the ID with gh api /orgs/<org>/hooks or gh api /repos/<owner>/<repo>/hooks"
      >
        <Input
          aria-label={components.QueryEditor.WebhookDeliveries.hookIdInput}
          width={RightColumnWidth}
          value={hookId}
          onChange={(el) => setHookId(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, hookId: el.currentTarget.value })}
        />
      </EditorField>
      <EditorField
        label="Event"
        tooltip="Only return deliveries triggered by this event (e.g. pull_request). Leave empty for all events"
      >
        <Input
          aria-label={components.QueryEditor.WebhookDeliveries.eventInput}
          width={RightColumnWidth}
          value={event}
          onChange={(el) => setEvent(el.currentTarget.value)}
          onBlur={(el) => props.onChange({ ...props, event: el.currentTarget.value })}
        />
      </EditorField>
      <EditorField
        label="Status"
        tooltip="Deliveries answered with a 2xx status code count as successful. Deliveries that never got a response count as failed"
      >
        <Combobox
          options={statusOptions}
          width={RightColumnWidth}
          value={props.status || WebhookDeliveryStatus.All}
          onChange={(opt) => props.onChange({ ...props, status: opt.value as WebhookDeliveryStatus })}
        />
      </EditorField>
    </EditorRow>
  );
};
