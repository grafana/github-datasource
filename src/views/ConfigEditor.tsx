import React, { PureComponent } from 'react';
import { onUpdateDatasourceJsonDataOption, DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { Input, InlineFormLabel, LegacyForms, InfoBox, Icon } from '@grafana/ui';
import { GithubDataSourceOptions, GithubSecureJsonData } from '../types';
import { selectors } from 'components/selectors';

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
        <InfoBox title="Access Token Permissions">
          <p>
            To create a new Access Token, navigate to{' '}
            <a href="https://github.com/settings/tokens">
              Personal Access Tokens <Icon name="link" />
            </a>{' '}
            and create a click &quot;Generate new token.&quot;
          </p>
          <p>Ensure that your token has the following permissions:</p>
          <h4>For all repositories:</h4>
          <pre>
            <ul>
              <li>public_repo</li>
              <li>repo:status</li>
              <li>repo_deployment</li>
              <li>read:packages</li>
            </ul>
            <ul>
              <li>user:read</li>
              <li>user:email</li>
            </ul>
          </pre>
          <h4>An extra setting is required for private repositories:</h4>
          <pre>
            <ul>
              <li>repo (Full control of private repositories)</li>
            </ul>
          </pre>
        </InfoBox>
        <div className="gf-form-group">
          <h3 className="page-heading">Service Account Access</h3>
          <div className="gf-form">
            <LegacyForms.SecretFormField
              aria-label={selectors.components.ConfigEditor.AccessToken.input}
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
          <h3 className="page-heading">GitHub Enterprise</h3>
          <div className="gf-form">
            <InlineFormLabel className="width-10">GitHub Enterprise URL</InlineFormLabel>
            <Input
              css=""
              className="width-27"
              value={jsonData.githubUrl}
              placeholder="URL of GitHub Enterprise"
              summary="URL for GitHub Enterprise, such as https://github.company.com, leave blank if using github.com"
              onChange={onUpdateDatasourceJsonDataOption(this.props, 'githubUrl')}
            />
          </div>
        </div>
        <InfoBox title="GitHub Enterprise">
          <p>For GitHub Enterprise enter the URL, such as https://github.company.com</p>
          <p>Leave blank if not using GitHub Enterprise, which will default to github.com</p>
        </InfoBox>
        {/*<div className="gf-form-group">
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
        </div>*/}
      </>
    );
  }
}

export default ConfigEditor;
