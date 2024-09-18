import { css } from '@emotion/css';
import {
  DataSourcePluginOptionsEditorProps,
  GrafanaTheme2,
  onUpdateDatasourceJsonDataOption,
  onUpdateDatasourceSecureJsonDataOption,
} from '@grafana/data';
import { ConfigSection, DataSourceDescription } from '@grafana/experimental';
import { Collapse, Field, Input, Label, RadioButtonGroup, SecretInput, SecretTextArea, useStyles2 } from '@grafana/ui';
import React, { ChangeEvent, useState } from 'react';
import { components } from '../components/selectors';
import { GithubDataSourceOptions, GithubSecureJsonData } from '../types';
import { Divider } from 'components/Divider';

export type ConfigEditorProps = DataSourcePluginOptionsEditorProps<GithubDataSourceOptions, GithubSecureJsonData>;

const ConfigEditor = (props: ConfigEditorProps) => {
  const { options, onOptionsChange } = props;
  const { jsonData, secureJsonData, secureJsonFields } = options;
  const secureSettings = (secureJsonData || {}) as GithubSecureJsonData;
  const styles = useStyles2(getStyles);
  const WIDTH_LONG = 40;
  const authOptions = [
    { label: 'Personal Access Token', value: 'personal-access-token' },
    { label: 'GitHub App', value: 'github-app' },
  ];
  const licenseOptions = [
    { label: 'Basic', value: 'github-basic' },
    { label: 'Enterprise', value: 'github-enterprise' },
  ];

  const [isOpen, setIsOpen] = useState(true);
  const [selectedLicense, setSelectedLicense] = useState(jsonData.githubUrl ? 'github-enterprise' : 'github-basic');
  const [selectedAuthType, setSelectedAuthType] = useState(
    secureJsonFields!['privateKey'] ? 'github-app' : 'personal-access-token'
  );

  const onSettingUpdate = (prop: string, set = true) => {
    return (event: ChangeEvent<HTMLInputElement>) => {
      const { onOptionsChange, options } = props;
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
  };

  const onSettingReset = (prop: string) => () => {
    onSettingUpdate(prop, false)({ target: { value: '' } } as ChangeEvent<HTMLInputElement>);
  };

  const onAuthChange = (value: string) => {
    // clear out personal access token when switching to app auth
    if (value === 'github-app') {
      onSettingReset('accessToken')();
    }

    // clear out app token when switching to personal access token
    if (value === 'personal-access-token') {
      onSettingReset('privateKey')();
    }

    setSelectedAuthType(value);
  };

  const onLicenseChange = (value: string) => {
    // clear out githubUrl when switching to basic
    if (value === 'github-basic') {
      onOptionsChange({ ...options, jsonData: { ...jsonData, githubUrl: '' } });
    }

    setSelectedLicense(value);
  };

  return (
    <>
      <DataSourceDescription
        dataSourceName="GitHub"
        docsLink="https://grafana.com/grafana/plugins/grafana-github-datasource/"
        hasRequiredFields={false}
      />

      <Divider />

      <Collapse collapsible label="Access Token Permissions" isOpen={isOpen} onToggle={() => setIsOpen((x) => !x)}>
        <p>
          To create a new Access Token, navigate to{' '}
          <a className={styles.externalLink} href="https://github.com/settings/tokens" target="_blank" rel="noreferrer">
            Personal Access Tokens
          </a>{' '}
          and create a click &quot;Generate new token.&quot;
        </p>

        <p>Ensure that your token has the following permissions:</p>

        <b>For all repositories:</b>
        <pre>public_repo, repo:status, repo_deployment, read:packages, read:user, user:email</pre>

        <b>For Github projects:</b>
        <pre>read:org, read:project</pre>

        <b>An extra setting is required for private repositories:</b>
        <pre>repo (Full control of private repositories)</pre>
      </Collapse>

      <Divider />

      <ConfigSection title="Connection">
        <RadioButtonGroup
          options={authOptions}
          value={selectedAuthType}
          onChange={onAuthChange}
          className={styles.radioButton}
        />

        {selectedAuthType === 'personal-access-token' && (
          <Field label="Personal Access Token">
            <SecretInput
              placeholder="Personal Access Token"
              data-testid={components.ConfigEditor.AccessToken}
              value={secureSettings.accessToken || ''}
              isConfigured={secureJsonFields!['accessToken']}
              onChange={onSettingUpdate('accessToken', false)}
              onReset={onSettingReset('accessToken')}
              width={WIDTH_LONG}
            />
          </Field>
        )}

        {selectedAuthType === 'github-app' && (
          <>
            <Field label="App ID">
              <Input
                placeholder="App ID"
                value={jsonData.appId}
                onChange={onUpdateDatasourceJsonDataOption(props, 'appId')}
                width={WIDTH_LONG}
              />
            </Field>
            <Field label="Installation ID">
              <Input
                placeholder="Installation ID"
                value={jsonData.installationId}
                onChange={onUpdateDatasourceJsonDataOption(props, 'installationId')}
                width={WIDTH_LONG}
              />
            </Field>
            <Field label="Private Key">
              <SecretTextArea
                placeholder="-----BEGIN CERTIFICATE-----"
                isConfigured={secureJsonFields!['privateKey']}
                onChange={onUpdateDatasourceSecureJsonDataOption(props, 'privateKey')}
                onReset={onSettingReset('privateKey')}
                cols={55}
                rows={7}
              />
            </Field>
          </>
        )}
      </ConfigSection>

      <Divider />

      <ConfigSection title="Additional Settings" isCollapsible>
        <Label>GitHub License</Label>
        <RadioButtonGroup
          options={licenseOptions}
          value={selectedLicense}
          onChange={onLicenseChange}
          className={styles.radioButton}
        />

        {selectedLicense === 'github-enterprise' && (
          <Field label="GitHub Enterprise URL">
            <Input
              placeholder="URL of GitHub Enterprise"
              value={jsonData.githubUrl}
              onChange={onUpdateDatasourceJsonDataOption(props, 'githubUrl')}
              width={WIDTH_LONG}
            />
          </Field>
        )}
      </ConfigSection>
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  return {
    externalLink: css`
      text-decoration: underline;
    `,
    radioButton: css`
      margin-bottom: ${theme.spacing(2)};
    `,
  };
};

export default ConfigEditor;
