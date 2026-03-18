import type { Spec } from '@json-render/react';
import configSchema from '../../pkg/schema/config.json';

/*
 * ---------------------------------------------------------------------------
 *  Config-editor UI spec – derived from the config JSON schema
 * ---------------------------------------------------------------------------
 *
 *  The generated JSON schema (pkg/schema/config.json) describes the DATA
 *  shape as a set of discriminated unions:
 *
 *   jsonData = allOf(
 *     anyOf(basic-plan | enterprise-cloud | enterprise-server),   ← plan
 *     anyOf(personal-access-token | github-app)                   ← auth
 *   )
 *
 *   secureJsonData = anyOf(
 *     { accessToken }   ← PAT
 *     | { privateKey }  ← GitHub App
 *   )
 *
 *  The function below walks the schema to extract field names, descriptions,
 *  enum constants and discriminator values so that the UI spec is driven by
 *  the schema rather than hard-coded labels.
 * ---------------------------------------------------------------------------
 */

// ---------------------------------------------------------------------------
// Helpers – extract metadata from the JSON schema
// ---------------------------------------------------------------------------

type JSONSchemaObject = Record<string, unknown>;
type JSONSchemaProperty = { type?: string; const?: string; description?: string; not?: unknown };

/** Safely narrow an `unknown` to a JSON-schema-like object. */
function isObj(v: unknown): v is JSONSchemaObject {
  return typeof v === 'object' && v !== null && !Array.isArray(v);
}

/** Return the `properties` record of a JSON-schema object node. */
function propsOf(node: JSONSchemaObject): Record<string, JSONSchemaProperty> {
  const p = node['properties'];
  return isObj(p) ? (p as Record<string, JSONSchemaProperty>) : {};
}

/**
 * Collect every concrete variant from the schema's anyOf/allOf structure
 * for a given top-level key (e.g. "jsonData" or "secureJsonData").
 */
function collectVariants(topKey: 'jsonData' | 'secureJsonData'): JSONSchemaObject[] {
  const top = (configSchema as JSONSchemaObject).properties;
  if (!isObj(top)) {
    return [];
  }
  const section = top[topKey];
  if (!isObj(section)) {
    return [];
  }

  const result: JSONSchemaObject[] = [];

  function walk(node: unknown) {
    if (!isObj(node)) {
      return;
    }
    if (Array.isArray(node['anyOf'])) {
      (node['anyOf'] as unknown[]).forEach(walk);
    }
    if (Array.isArray(node['allOf'])) {
      (node['allOf'] as unknown[]).forEach(walk);
    }
    // A leaf variant has `properties`
    if (isObj(node['properties'])) {
      result.push(node as JSONSchemaObject);
    }
  }

  walk(section);
  return result;
}

// ---------------------------------------------------------------------------
// Extract auth-type & plan-type options from the schema
// ---------------------------------------------------------------------------

const jsonDataVariants = collectVariants('jsonData');
const secureVariants = collectVariants('secureJsonData');

/** Return the authentication-type variants (discriminator: selectedAuthType). */
function authVariants() {
  return jsonDataVariants.filter((v) => 'selectedAuthType' in propsOf(v));
}

/** Return the plan/license variants (discriminator: githubPlan). */
function planVariants() {
  return jsonDataVariants.filter((v) => 'githubPlan' in propsOf(v));
}

/**
 * Build radio-group options from an array of schema variants keyed by a
 * discriminator property.
 */
function radioOptionsFromVariants(
  variants: JSONSchemaObject[],
  discriminator: string,
  labelMap: Record<string, string>
): Array<{ label: string; value: string }> {
  return variants
    .map((v) => {
      const prop = propsOf(v)[discriminator];
      const constVal = prop?.const;
      if (!constVal) {
        return null;
      }
      return { label: labelMap[constVal] ?? constVal, value: constVal };
    })
    .filter(Boolean) as Array<{ label: string; value: string }>;
}

// Human-readable labels keyed by schema const values
const AUTH_LABELS: Record<string, string> = {
  'personal-access-token': 'Personal Access Token',
  'github-app': 'GitHub App',
};

const PLAN_LABELS: Record<string, string> = {
  'github-basic': 'Free, Pro & Team',
  'github-enterprise-cloud': 'Enterprise Cloud',
  'github-enterprise-server': 'Enterprise Server',
};

// ---------------------------------------------------------------------------
// Build the json-render spec
// ---------------------------------------------------------------------------

/**
 * Build the complete json-render `Spec` for the config editor.
 *
 * State keys map 1-to-1 to jsonData / secureJsonData property names so
 * that the bridge in `ConfigEditorJsonRender.tsx` can propagate changes
 * without a mapping layer.
 */
export function buildConfigEditorSpec(state: {
  selectedAuthType: string;
  githubPlan: string;
  githubUrl: string;
  appId: string;
  installationId: string;
  accessTokenConfigured: boolean;
  privateKeyConfigured: boolean;
}): Spec {
  // Derive descriptions from schema for field labels / placeholders
  const authVariantsArr = authVariants();
  const planVariantsArr = planVariants();

  const ghAppVariant = authVariantsArr.find((v) => propsOf(v)['selectedAuthType']?.const === 'github-app');
  const enterpriseServerVariant = planVariantsArr.find(
    (v) => propsOf(v)['githubPlan']?.const === 'github-enterprise-server'
  );

  const appIdDesc = ghAppVariant ? propsOf(ghAppVariant)['appId']?.description ?? 'App ID' : 'App ID';
  const installationIdDesc = ghAppVariant
    ? propsOf(ghAppVariant)['installationId']?.description ?? 'Installation ID'
    : 'Installation ID';
  const githubUrlDesc = enterpriseServerVariant
    ? propsOf(enterpriseServerVariant)['githubUrl']?.description ?? 'GitHub Enterprise Server URL'
    : 'GitHub Enterprise Server URL';

  // Schema-derived descriptions for secure fields
  const patSecureVariant = secureVariants.find((v) => {
    const p = propsOf(v);
    return p['accessToken'] && !p['accessToken'].not;
  });
  const ghAppSecureVariant = secureVariants.find((v) => {
    const p = propsOf(v);
    return p['privateKey'] && !p['privateKey'].not;
  });
  const accessTokenDesc = patSecureVariant
    ? propsOf(patSecureVariant)['accessToken']?.description ?? 'Personal Access Token'
    : 'Personal Access Token';
  const privateKeyDesc = ghAppSecureVariant
    ? propsOf(ghAppSecureVariant)['privateKey']?.description ?? 'Private Key'
    : 'Private Key';

  const authOptions = radioOptionsFromVariants(authVariantsArr, 'selectedAuthType', AUTH_LABELS);
  const planOptions = radioOptionsFromVariants(planVariantsArr, 'githubPlan', PLAN_LABELS);

  return {
    root: 'root',
    state: {
      selectedAuthType: state.selectedAuthType,
      githubPlan: state.githubPlan,
      githubUrl: state.githubUrl,
      appId: state.appId,
      installationId: state.installationId,
      accessTokenConfigured: state.accessTokenConfigured,
      privateKeyConfigured: state.privateKeyConfigured,
    },
    elements: {
      root: {
        type: 'Stack',
        props: {},
        children: ['description', 'divider-1', 'auth-section', 'divider-2', 'connection-section'],
      },

      // -- Header --
      description: {
        type: 'Description',
        props: {
          dataSourceName: 'GitHub',
          docsLink: 'https://grafana.com/docs/plugins/grafana-github-datasource',
        },
        children: [],
      },

      'divider-1': { type: 'Divider', props: {}, children: [] },

      // -- Authentication (driven by jsonData auth anyOf) --
      'auth-section': {
        type: 'Section',
        props: { title: 'Authentication' },
        children: ['auth-type', 'pat-token', 'app-id', 'installation-id', 'private-key'],
      },

      'auth-type': {
        type: 'RadioGroup',
        props: {
          label: 'Authentication Type',
          value: { $bindState: '/selectedAuthType' } as unknown as string,
          options: authOptions,
        },
        children: [],
      },

      // PAT fields – visible when selectedAuthType == "personal-access-token"
      'pat-token': {
        type: 'SecretInput',
        props: {
          label: accessTokenDesc,
          placeholder: 'Personal Access Token',
          isConfigured: { $state: '/accessTokenConfigured' } as unknown as boolean,
        },
        visible: { $state: '/selectedAuthType', eq: 'personal-access-token' },
        children: [],
      },

      // GitHub App fields – visible when selectedAuthType == "github-app"
      'app-id': {
        type: 'TextInput',
        props: {
          label: appIdDesc,
          value: { $bindState: '/appId' } as unknown as string,
          placeholder: 'App ID',
        },
        visible: { $state: '/selectedAuthType', eq: 'github-app' },
        children: [],
      },

      'installation-id': {
        type: 'TextInput',
        props: {
          label: installationIdDesc,
          value: { $bindState: '/installationId' } as unknown as string,
          placeholder: 'Installation ID',
        },
        visible: { $state: '/selectedAuthType', eq: 'github-app' },
        children: [],
      },

      'private-key': {
        type: 'SecretTextArea',
        props: {
          label: privateKeyDesc,
          placeholder: '-----BEGIN CERTIFICATE-----',
          isConfigured: { $state: '/privateKeyConfigured' } as unknown as boolean,
        },
        visible: { $state: '/selectedAuthType', eq: 'github-app' },
        children: [],
      },

      'divider-2': { type: 'Divider', props: {}, children: [] },

      // -- Connection / License (driven by jsonData plan anyOf) --
      'connection-section': {
        type: 'Section',
        props: { title: 'Connection', isCollapsible: true },
        children: ['license-type', 'enterprise-url'],
      },

      'license-type': {
        type: 'RadioGroup',
        props: {
          label: 'GitHub License Type',
          value: { $bindState: '/githubPlan' } as unknown as string,
          options: planOptions,
        },
        children: [],
      },

      'enterprise-url': {
        type: 'TextInput',
        props: {
          label: githubUrlDesc,
          value: { $bindState: '/githubUrl' } as unknown as string,
          placeholder: 'http(s)://HOSTNAME/',
        },
        visible: { $state: '/githubPlan', eq: 'github-enterprise-server' },
        children: [],
      },
    },
  };
}
