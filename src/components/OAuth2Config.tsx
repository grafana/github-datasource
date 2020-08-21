import React, { useMemo, useState } from 'react';
import { Input, Button } from '@grafana/ui';
import { QueryInlineField } from './Forms';
import { ConfigEditorProps } from '../types';
import { onUpdateDatasourceSecureJsonDataOption, onUpdateDatasourceResetOption } from '@grafana/data';
import ResetButton from './ResetButton';

const copyToClipboard = async (text: string) => {
  navigator.clipboard.writeText(text).then(
    function() {
      // TODO: toast, successful copy
    },
    function(err) {
      console.error(err);
    }
  );
};

export default (props: ConfigEditorProps) => {
  const clientIDSet = useMemo(() => {
    return props.options.secureJsonFields?.clientID;
  }, [props.options.secureJsonFields?.clientID]);

  const clientSecretSet = useMemo(() => {
    return props.options.secureJsonFields?.clientSecret;
  }, [props.options.secureJsonFields?.clientSecret]);

  const accessTokenSet = useMemo(() => {
    return props.options.secureJsonFields?.oauthAccessToken;
  }, [props.options.secureJsonFields?.oauthAccessToken]);

  const [rootURL, setRootURL] = useState<string>(`${window.location.protocol}//${window.location.host}`);

  const authCallbackURL = useMemo(() => {
    return `${rootURL}/api/datasources/1631/resources/auth/callback`;
  }, [rootURL]);

  return (
    <>
      <QueryInlineField
        label="Grafana Root URL"
        labelWidth={12}
        tooltip="The root URL of your Grafana instance used for the Authorization callback"
      >
        <Input
          css=""
          width={56}
          value={rootURL}
          placeholder="https://grafana.com"
          onChange={e => setRootURL(e.currentTarget.value)}
        />
      </QueryInlineField>
      <QueryInlineField label="Authorization Callback" labelWidth={12}>
        <Input css="" width={56} disabled value={authCallbackURL} />
        <Button type="button" variant="secondary" onClick={() => copyToClipboard(authCallbackURL)}>
          Copy
        </Button>
      </QueryInlineField>
      <QueryInlineField label="Client ID" labelWidth={12}>
        <Input
          css=""
          disabled={clientIDSet}
          width={32}
          placeholder={clientIDSet ? 'Stored securely' : 'Client ID'}
          onChange={onUpdateDatasourceSecureJsonDataOption(props, 'clientID')}
          value={props.options.secureJsonData?.clientID}
        />
        <ResetButton disabled={!clientIDSet} onClick={onUpdateDatasourceResetOption(props, 'clientID')} />
      </QueryInlineField>
      <QueryInlineField label="Client Secret" labelWidth={12}>
        <Input
          css=""
          disabled={clientSecretSet}
          width={32}
          placeholder={clientSecretSet ? 'Stored securely' : 'Client Secret'}
          onChange={onUpdateDatasourceSecureJsonDataOption(props, 'clientSecret')}
          value={props.options.secureJsonData?.clientSecret}
        />
        <ResetButton disabled={!clientSecretSet} onClick={onUpdateDatasourceResetOption(props, 'clientSecret')} />
      </QueryInlineField>
      <QueryInlineField label="Access Token" labelWidth={12} tooltip='To populate this field, click "Authorize"'>
        <Input
          css=""
          disabled={true}
          width={32}
          placeholder={accessTokenSet ? 'Stored securely' : 'Access token...'}
          onChange={onUpdateDatasourceSecureJsonDataOption(props, 'clientSecret')}
          value={props.options.secureJsonData?.oauthAccessToken}
        />
        <ResetButton disabled={!accessTokenSet} onClick={onUpdateDatasourceResetOption(props, 'oauthAccessToken')} />
      </QueryInlineField>
      <QueryInlineField
        label="Authorize Github"
        labelWidth={12}
        tooltip='Authorizing to GitHub using OAuth 2.0 requires "Client ID" and "Client Secret" to be set.'
      >
        <Button type="button" variant="primary" disabled={!(clientIDSet && clientSecretSet)}>
          Authorize
        </Button>
      </QueryInlineField>
    </>
  );
};
