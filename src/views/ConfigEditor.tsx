import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceJsonDataOption } from '@grafana/data';
import { InlineFormLabel, Input, LegacyForms } from '@grafana/ui';
import { GithubDataSourceOptions, GithubSecureJsonData } from '../types';

export type ConfigEditorProps = DataSourcePluginOptionsEditorProps<GithubDataSourceOptions, GithubSecureJsonData>;

export class ConfigEditor extends PureComponent<ConfigEditorProps> {
  onSettingReset = (prop: string) => (event: any) => {
    this.onSettingUpdate(prop, false)({ target: { value: undefined } });
  };

  onSettingUpdate = (prop: string, set = true) => (event: any) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        [prop]: event.target.value,
      },
      secureJsonFields: {
        ...options.secureJsonFields,
        [prop]: set,
      },
    });
  };

  render() {
    const {
      options: { jsonData, secureJsonData, secureJsonFields },
    } = this.props;
    const secureSettings = (secureJsonData || {}) as GithubSecureJsonData;
    return (
      <>
        <div className="gf-form-group">
          <h3 className="page-heading">Service Account Access</h3>
          {/* <div className="gf-form">
          <LegacyForms.FormField
            label="API URL"
            labelWidth={11}
            inputWidth={27}
            tooltip={'URL to Datadog API'}
            onChange={onUpdateDatasourceJsonDataOption(this.props, 'url')}
            value={jsonData.url || 'https://api.datadoghq.com'}
            placeholder="https://api.datadoghq.com"
          />
        </div> */}
          <div className="gf-form">
            <LegacyForms.SecretFormField
              label="Access Token"
              inputWidth={27}
              labelWidth={10}
              onChange={this.onSettingUpdate('accessToken', false)}
              onBlur={this.onSettingUpdate('accessToken')}
              value={secureSettings.accessToken || ''}
              placeholder="Github Personal Access Token"
              onReset={this.onSettingReset('accessToken')}
              isConfigured={secureJsonFields!['accessToken']}
            />
          </div>
        </div>

        <div className="gf-form-group">
          <h3 className="page-heading">Default Query Options</h3>
          <div className="gf-form">
            <InlineFormLabel className="width-10">Owner</InlineFormLabel>
            <Input
              css=""
              className="width-9"
              value={jsonData.owner}
              placeholder="username or organization"
              onChange={onUpdateDatasourceJsonDataOption(this.props, 'owner')}
            />
          </div>
          <div className="gf-form">
            <InlineFormLabel className="width-10">Repository</InlineFormLabel>
            <Input
              css=""
              className="width-9"
              value={jsonData.repository}
              placeholder="the repo name"
              onChange={onUpdateDatasourceJsonDataOption(this.props, 'repository')}
            />
          </div>
        </div>
      </>
    );
  }
}

export default ConfigEditor;
