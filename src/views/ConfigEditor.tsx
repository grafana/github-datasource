import { css } from '@emotion/css';
import { DataSourcePluginOptionsEditorProps, GrafanaTheme2, onUpdateDatasourceJsonDataOption } from '@grafana/data';
import { ConfigSection, DataSourceDescription } from '@grafana/experimental';
import { Collapse, Divider, Field, Input, Label, RadioButtonGroup, SecretInput, useStyles2 } from '@grafana/ui';
import React, { ChangeEvent, useState } from 'react';
import { selectors } from '../components/selectors';
import { GithubDataSourceOptions, GithubSecureJsonData } from '../types';

export type ConfigEditorProps = DataSourcePluginOptionsEditorProps<GithubDataSourceOptions, GithubSecureJsonData>;

const ConfigEditor = (props: ConfigEditorProps) => {
  const { jsonData, secureJsonData, secureJsonFields } = props.options;
  const secureSettings = (secureJsonData || {}) as GithubSecureJsonData;
  const styles = useStyles2(getStyles);

  const [selectedLicense, setSelectedLicense] = useState(jsonData.githubUrl ? 'github-enterprise' : 'github-basic');
  const [isOpen, setIsOpen] = useState(false);

  const onSettingUpdate =
    (prop: string, set = true) =>
    (event: ChangeEvent<HTMLInputElement>) => {
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

  const onSettingReset = (prop: string) => () => {
    onSettingUpdate(prop, false)({ target: { value: '' } } as ChangeEvent<HTMLInputElement>);
  };

  const onLicenseChange = (value: string) => {
    if (value === 'github-basic') {
      jsonData.githubUrl = '';
    }

    setSelectedLicense(value);
  };

  const licenseOptions = [
    { label: 'Basic', value: 'github-basic' },
    { label: 'Enterprise', value: 'github-enterprise' },
  ];

  const WIDTH_LONG = 40;

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
        <Field label="Access Token">
          <SecretInput
            placeholder="GitHub Personal Access Token"
            data-testid={selectors.components.ConfigEditor.AccessToken}
            value={secureSettings.accessToken || ''}
            isConfigured={secureJsonFields!['accessToken']}
            onChange={onSettingUpdate('accessToken', false)}
            onReset={onSettingReset('accessToken')}
            width={WIDTH_LONG}
          />
        </Field>
      </ConfigSection>

      <Divider />

      <ConfigSection title="Additional Settings">
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
