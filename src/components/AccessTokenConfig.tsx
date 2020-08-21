import React, { useMemo } from 'react';
import { QueryInlineField } from './Forms';
import ResetButton from './ResetButton';
import { ConfigEditorProps } from '../types';

import { Input } from '@grafana/ui';
import { onUpdateDatasourceSecureJsonDataOption, onUpdateDatasourceResetOption } from '@grafana/data';

export default (props: ConfigEditorProps) => {
  const accessTokenSet = useMemo(() => {
    return props.options.secureJsonFields?.accessToken;
  }, [props.options.secureJsonFields?.accessToken]);

  return (
    <>
      <QueryInlineField label="Access Token" labelWidth={12}>
        <Input
          placeholder="personal access token..."
          css=""
          width={48}
          disabled={accessTokenSet}
          value={props.options.secureJsonData?.accessToken}
          onChange={onUpdateDatasourceSecureJsonDataOption(props, 'accessToken')}
          type="password"
        />
        <ResetButton disabled={!accessTokenSet} onClick={onUpdateDatasourceResetOption(props, 'accessToken')} />
      </QueryInlineField>
    </>
  );
};
