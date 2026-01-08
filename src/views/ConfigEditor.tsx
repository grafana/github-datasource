import React, { ChangeEvent, useCallback, useEffect, useState } from 'react';
import { css } from '@emotion/css';
import {
  onUpdateDatasourceJsonDataOption,
  onUpdateDatasourceSecureJsonDataOption,
  type DataSourcePluginOptionsEditorProps,
  type GrafanaTheme2,
  type SelectableValue,
} from '@grafana/data';
import { ConfigSection, DataSourceDescription } from '@grafana/plugin-ui';
import { config } from '@grafana/runtime';
import {
  Collapse,
  Field,
  Input,
  Label,
  RadioButtonGroup,
  SecretInput,
  SecretTextArea,
  SecureSocksProxySettings,
  useStyles2,
} from '@grafana/ui';
import { Divider } from 'components/Divider';
import type { GitHubAuthType, GitHubLicenseType, GitHubDataSourceOptions, GitHubSecureJsonData } from 'types/config';
import { components as selectors } from '../components/selectors';

export type ConfigEditorProps = DataSourcePluginOptionsEditorProps<GitHubDataSourceOptions, GitHubSecureJsonData>;

const ConfigEditor = (props: ConfigEditorProps) => {
  const { options, onOptionsChange } = props;
  const { jsonData, secureJsonData, secureJsonFields } = options;
  const secureSettings = secureJsonData || {};
  const styles = useStyles2(getStyles);
  const WIDTH_LONG = 40;

  const authOptions: Array<SelectableValue<GitHubAuthType>> = [
    { label: 'Personal Access Token', value: 'personal-access-token' },
    { label: 'GitHub App', value: 'github-app' },
  ];

  const licenseOptions: Array<SelectableValue<GitHubLicenseType>> = [
    { label: 'Free, Pro & Team', value: 'github-basic' },
    { label: 'Enterprise Cloud', value: 'github-enterprise-cloud' },
    { label: 'Enterprise Server', value: 'github-enterprise-server' },
  ];

  const [isOpen, setIsOpen] = useState(true);

  // Previously we used only githubUrl property to determine if the github plan is enterprise which is incorrect way
  // Also only on prem github enterprise will be having their own base URLs where as cloud will be having common URL and this causes confusions to the user
  // So we are adding a new prop called githubPlan to determine if the github instance is on-prem / cloud / basic plan
  // Also if no plan exist and no url exist, we need to fallback to github-basic
  // https://docs.github.com/en/get-started/using-github-docs/about-versions-of-github-docs
  const [selectedLicense, setSelectedLicense] = useState<GitHubLicenseType>(
    jsonData.githubPlan === 'github-enterprise-server' || jsonData.githubUrl
      ? 'github-enterprise-server'
      : jsonData?.githubPlan || 'github-basic'
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

  const onAuthChange = useCallback(
    (value: GitHubAuthType) => {
      onOptionsChange({ ...options, jsonData: { ...jsonData, selectedAuthType: value } });
    },
    [jsonData, onOptionsChange, options]
  );

  const onLicenseChange = (githubPlan: GitHubLicenseType) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        githubPlan,
        githubUrl: githubPlan === 'github-enterprise-server' ? jsonData.githubUrl : '',
      },
    });
    setSelectedLicense(githubPlan);
  };

  useEffect(() => {
    // set the default auth type if its a new datasource and nothing is set
    if (!jsonData.selectedAuthType) {
      onAuthChange('personal-access-token');
    }
  }, [jsonData.selectedAuthType, onAuthChange]);

  return (
    <>
      <DataSourceDescription
        dataSourceName="GitHub"
        docsLink="https://grafana.com/docs/plugins/grafana-github-datasource"
        hasRequiredFields={false}
      />

      <Divider />

      <Collapse collapsible label="Access Token & Permissions" isOpen={isOpen} onToggle={() => setIsOpen((x) => !x)}>
        <h4>How to create a access token</h4>
        <p>
          To create a new fine grained access token, navigate to{' '}
          <a
            className={styles.externalLink}
            href="https://github.com/settings/personal-access-tokens/new"
            target="_blank"
            rel="noreferrer"
          >
            Personal Access Tokens
          </a>{' '}
          or refer the guidelines from{' '}
          <a
            className={styles.externalLink}
            href="https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token"
            target="_blank"
            rel="noreferrer"
          >
            the Github documentation.
          </a>
        </p>
        <h4>Repository access</h4>
        <p>
          In the <b>Repository access</b> section, Select the required repositories you want to use with the plugin.
        </p>
        <h4>Permissions</h4>
        <p>
          In the repository permissions, Ensure to provide <b>read-only access</b> to the necessary section which you
          want to use with the plugin. <b>The plugin does not require any write access.</b> <br />
          Along with other permissions such as `Issues`, `Pull Requests`, ensure to provide read-only access to `Meta
          data` section as well.
          <br />
          This plugin does not require any org level permissions
        </p>
      </Collapse>
      <Divider />

      <ConfigSection title="Authentication">
        <Field label="Authentication Type">
          <RadioButtonGroup<GitHubAuthType>
            options={authOptions}
            value={jsonData.selectedAuthType}
            onChange={onAuthChange}
            className={styles.radioButton}
          />
        </Field>

        {jsonData.selectedAuthType === 'personal-access-token' && (
          <Field label="Personal Access Token">
            <SecretInput
              placeholder="Personal Access Token"
              data-testid={selectors.ConfigEditor.AccessToken}
              value={secureSettings.accessToken || ''}
              isConfigured={secureJsonFields!['accessToken']}
              onChange={onSettingUpdate('accessToken', false)}
              onReset={onSettingReset('accessToken')}
              width={WIDTH_LONG}
            />
          </Field>
        )}

        {jsonData.selectedAuthType === 'github-app' && (
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
      {config.secureSocksDSProxyEnabled && (
        <SecureSocksProxySettings options={options} onOptionsChange={onOptionsChange} />
      )}
      <Divider />

      <ConfigSection title="Connection" isCollapsible>
        <Label>GitHub License Type</Label>
        <RadioButtonGroup
          options={licenseOptions}
          value={selectedLicense}
          onChange={onLicenseChange}
          className={styles.radioButton}
        />

        {selectedLicense === 'github-enterprise-server' && (
          <Field label="GitHub Enterprise Server URL">
            <Input
              placeholder="http(s)://HOSTNAME/"
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
