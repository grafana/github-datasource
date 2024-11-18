import React, { ChangeEvent, useCallback, useEffect, useState } from 'react';
import { css } from '@emotion/css';
import { Collapse, Field, Input, Label, RadioButtonGroup, SecretInput, SecretTextArea, useStyles2 } from '@grafana/ui';
import { ConfigSection, DataSourceDescription } from '@grafana/experimental';
import { Divider } from 'components/Divider';
import { components as selectors } from '../components/selectors';
import {
  onUpdateDatasourceJsonDataOption,
  onUpdateDatasourceSecureJsonDataOption,
  type DataSourcePluginOptionsEditorProps,
  type GrafanaTheme2,
  type SelectableValue,
} from '@grafana/data';
import type { GitHubAuthType, GitHubLicenseType, GitHubDataSourceOptions, GitHubSecureJsonData } from 'types/config';

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

  const [selectedLicense, setSelectedLicense] = useState<GitHubLicenseType>(
    jsonData.githubPlan || jsonData.githubUrl ? 'github-enterprise-server' : jsonData?.githubPlan || 'github-basic'
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

      <Collapse collapsible label="Access Token Permissions" isOpen={isOpen} onToggle={() => setIsOpen((x) => !x)}>
        <p>
          To create a new Access Token, navigate to{' '}
          <a
            className={styles.externalLink}
            href="https://github.com/settings/tokens?type=beta"
            target="_blank"
            rel="noreferrer"
          >
            Personal Access Tokens
          </a>{' '}
          and create a click &quot;Generate new token.&quot;
        </p>

        <p>Ensure that your token has the following permissions:</p>

        <b>For all repositories:</b>
        <pre>public_repo, repo:status, repo_deployment, read:packages, read:user, user:email</pre>

        <b>For GitHub projects:</b>
        <pre>read:org, read:project</pre>

        <b>An extra setting is required for private repositories:</b>
        <pre>repo (Full control of private repositories)</pre>
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
