import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { Input, Button } from '@grafana/ui';
import { DataSourceOptions, SecureJsonData } from '../types';
import { QueryInlineField } from './Forms';
interface Props extends DataSourcePluginOptionsEditorProps<DataSourceOptions, SecureJsonData> {};

export default class extends PureComponent<Props> {
  render() {
    return (
      <>
        <QueryInlineField label="Authorization Callback" labelWidth={12}>
          <Input css='' width={56} disabled value="http://localhost:3000/api/datasources/1631/resources/auth/callback" />
          <Button type='button' variant='secondary' onClick={() => {}}>Copy</Button>
        </QueryInlineField>
        <QueryInlineField label="Client ID" labelWidth={12}>
          <Input css='' width={32} />
        </QueryInlineField>
        <QueryInlineField label="Client Secret" labelWidth={12}>
          <Input css='' width={32} />
        </QueryInlineField>
        <QueryInlineField label="Authorize Github" labelWidth={12}>
          <Button type='button' variant='primary' disabled>Authorize</Button>
        </QueryInlineField>
      </>
    )
  }
}
