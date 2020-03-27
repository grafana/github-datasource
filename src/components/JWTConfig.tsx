import React, { useState } from 'react';
import { startCase, isObject } from 'lodash';
import { FormLabel, Button } from '@grafana/ui';
import { DropZone } from './';
// import { JWTFile } from '../types';

const configKeys = [
  'type',
  'project_id',
  'private_key_id',
  'private_key',
  'client_email',
  'client_id',
  'auth_uri',
  'token_uri',
  'auth_provider_x509_cert_url',
  'client_x509_cert_url',
];

export interface Props {
  onChange: (jwt: string) => void;
  isConfigured: boolean;
}

const validateJson = (json: { [key: string]: string }) => isObject(json) && configKeys.every(key => !!json[key]);

export function JWTConfig({ onChange, isConfigured }: Props) {
  const [enableUpload, setEnableUpload] = useState<boolean>(!isConfigured);
  const [error, setError] = useState<string>();

  return enableUpload ? (
    <>
      <DropZone
        baseStyle={{ marginTop: '24px' }}
        accept="application/json"
        onDrop={acceptedFiles => {
          const reader = new FileReader();
          if (acceptedFiles.length === 1) {
            reader.onloadend = (e: any) => {
              const json = JSON.parse(e.target.result);
              if (validateJson(json)) {
                onChange(e.target.result);
                setEnableUpload(false);
              } else {
                setError('Invalid JWT file');
              }
            };
            reader.readAsText(acceptedFiles[0]);
          } else if (acceptedFiles.length > 1) {
            setError('You can only upload one file');
          }
        }}
      >
        <p style={{ margin: 0, fontSize: 18 }}>Drop the file here, or click to use the file explorer</p>
      </DropZone>

      {error && (
        <pre style={{ margin: '12px 0 0' }} className="gf-form-pre alert alert-error">
          {error}
        </pre>
      )}
    </>
  ) : (
    <>
      {configKeys.map(key => (
        <div className="gf-form">
          <FormLabel width={10}>{startCase(key)}</FormLabel>
          <input disabled className="gf-form-input width-30" value="configured" />
        </div>
      ))}
      <Button style={{ marginTop: 12 }} variant="secondary" onClick={() => setEnableUpload(true)}>
        Upload another JWT file
      </Button>
    </>
  );
}
