import { DataSourceJsonData, DataSourcePluginOptionsEditorProps } from '@grafana/data';

export interface Label {
  color: string;
  description: string;
  name: string;
}

export interface DataSourceOptions {}

export interface DataSourceOptions extends DataSourceJsonData {}

export interface AuthConfig {
  // rootURL is used for generating the Authorization callback URL
  rootURL?: string;
  // accessToken is set if the user is using a Personal Access Token to connect to GitHub
  accessToken?: string;
  // oauthAccessToken is only set if the authorization flow is done through a GitHub OAuth 2.0 app
  oauthAccessToken?: string;
  // clientID is provided by GitHub when setting up a GitHub OAuth 2.0 app
  clientID?: string;
  // clientSecret is provided by GitHub when setting up a GitHub OAuth 2.0 app
  clientSecret?: string;
}

export interface SecureJsonData extends AuthConfig {}

export type ConfigEditorProps = DataSourcePluginOptionsEditorProps<DataSourceOptions, SecureJsonData>;
