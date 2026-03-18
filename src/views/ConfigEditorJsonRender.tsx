import React, { useCallback, useEffect, useMemo, useRef } from 'react';
import { Renderer, JSONUIProvider, createStateStore } from '@json-render/react';
import type { DataSourcePluginOptionsEditorProps, DataSourceSettings } from '@grafana/data';
import type { GitHubDataSourceOptions, GitHubSecureJsonData } from 'types/config';
import { createConfigEditorRegistry } from 'jsonrender/registry';
import { buildConfigEditorSpec } from 'jsonrender/configEditorSpec';

export type ConfigEditorProps = DataSourcePluginOptionsEditorProps<GitHubDataSourceOptions, GitHubSecureJsonData>;

/**
 * Proof-of-concept ConfigEditor rendered entirely via json-render.
 *
 * The component:
 *  1. Derives a json-render UI spec from the config JSON schema
 *     (see `configEditorSpec.ts`)
 *  2. Maps catalog components to Grafana UI via a component registry
 *     (see `registry.tsx`)
 *  3. Bridges json-render's internal state with Grafana's
 *     `onOptionsChange` callback so that edits persist.
 */
const ConfigEditorJsonRender = (props: ConfigEditorProps) => {
  const { options, onOptionsChange } = props;
  const { jsonData, secureJsonFields } = options;

  // ---- Stable references to avoid stale closures in callbacks ----
  const optionsRef = useRef(options);
  const onOptionsChangeRef = useRef(onOptionsChange);
  useEffect(() => {
    optionsRef.current = options;
    onOptionsChangeRef.current = onOptionsChange;
  }, [options, onOptionsChange]);

  // ---- Callbacks forwarded into the json-render action / registry layer ----

  const onJsonDataChange = useCallback((field: string, value: string) => {
    const opts = optionsRef.current;
    const newJsonData = { ...opts.jsonData, [field]: value } as GitHubDataSourceOptions;

    // When switching license, clear githubUrl for non-enterprise-server plans
    if (field === 'githubPlan' && value !== 'github-enterprise-server') {
      (newJsonData as unknown as Record<string, unknown>)['githubUrl'] = '';
    }

    onOptionsChangeRef.current({ ...opts, jsonData: newJsonData });
  }, []);

  const onSecureJsonDataChange = useCallback((field: string, value: string) => {
    const opts = optionsRef.current;
    onOptionsChangeRef.current({
      ...opts,
      secureJsonData: { ...opts.secureJsonData, [field]: value } as GitHubSecureJsonData,
      secureJsonFields: { ...opts.secureJsonFields, [field]: false },
    });
  }, []);

  const onResetSecureField = useCallback((field: string) => {
    const opts = optionsRef.current;
    onOptionsChangeRef.current({
      ...opts,
      secureJsonData: { ...opts.secureJsonData, [field]: '' } as GitHubSecureJsonData,
      secureJsonFields: { ...opts.secureJsonFields, [field]: false },
    });
  }, []);

  // ---- Registry (stable across renders because callbacks are ref-based) ----

  const { registry } = useMemo(
    () => createConfigEditorRegistry(onJsonDataChange, onSecureJsonDataChange, onResetSecureField),
    [onJsonDataChange, onSecureJsonDataChange, onResetSecureField]
  );

  // ---- State store (controlled mode) ----

  const store = useMemo(() => {
    return createStateStore(buildStateFromOptions(options));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Keep the store in sync with external option changes
  const prevOptionsRef = useRef(options);
  useEffect(() => {
    if (prevOptionsRef.current !== options) {
      prevOptionsRef.current = options;
      store.update(buildStateFromOptions(options));
    }
  }, [options, store]);

  // ---- Build spec (recalculated when relevant state changes) ----

  const selectedAuthType = (jsonData.selectedAuthType as string) || 'personal-access-token';
  const githubPlan = (jsonData.githubPlan as string) || 'github-basic';

  const jd = jsonData as unknown as Record<string, string>;
  const spec = useMemo(
    () =>
      buildConfigEditorSpec({
        selectedAuthType,
        githubPlan,
        githubUrl: jd.githubUrl ?? '',
        appId: jd.appId ?? '',
        installationId: jd.installationId ?? '',
        accessTokenConfigured: !!secureJsonFields?.['accessToken'],
        privateKeyConfigured: !!secureJsonFields?.['privateKey'],
      }),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [selectedAuthType, githubPlan, jsonData, secureJsonFields]
  );

  // ---- Set default auth type on first mount (matches original behaviour) ----

  useEffect(() => {
    if (!jsonData.selectedAuthType) {
      onJsonDataChange('selectedAuthType', 'personal-access-token');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <JSONUIProvider registry={registry} store={store}>
      <Renderer spec={spec} registry={registry} />
    </JSONUIProvider>
  );
};

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function buildStateFromOptions(
  options: DataSourceSettings<GitHubDataSourceOptions, GitHubSecureJsonData>
): Record<string, unknown> {
  const jd = options.jsonData as unknown as Record<string, unknown>;
  return {
    selectedAuthType: jd.selectedAuthType ?? 'personal-access-token',
    githubPlan: jd.githubPlan ?? 'github-basic',
    githubUrl: jd.githubUrl ?? '',
    appId: jd.appId ?? '',
    installationId: jd.installationId ?? '',
    accessTokenConfigured: !!options.secureJsonFields?.['accessToken'],
    privateKeyConfigured: !!options.secureJsonFields?.['privateKey'],
  };
}

export default ConfigEditorJsonRender;
