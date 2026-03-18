import { defineCatalog } from '@json-render/core';
import { schema } from '@json-render/react';
import { z } from 'zod';

/**
 * Component catalog for the config editor.
 *
 * Each component mirrors a Grafana UI primitive and exposes a Zod-typed
 * props contract so that json-render can validate specs at build-time and
 * generate AI prompts.
 */
export const configEditorCatalog = defineCatalog(schema, {
  components: {
    /** Root vertical stack container. */
    Stack: {
      props: z.object({}),
      description: 'A vertical stack that renders its children sequentially',
    },
    /** Titled collapsible section (maps to Grafana ConfigSection). */
    Section: {
      props: z.object({
        title: z.string(),
        isCollapsible: z.boolean().optional(),
      }),
      description: 'A titled section for grouping related form fields',
    },
    /** Radio button group for selecting one of several options. */
    RadioGroup: {
      props: z.object({
        label: z.string().optional(),
        value: z.string().optional(),
        options: z.array(z.object({ label: z.string(), value: z.string() })),
      }),
      description: 'A radio button group for selecting one of several predefined options',
    },
    /** Standard text input wrapped in a form field. */
    TextInput: {
      props: z.object({
        label: z.string(),
        value: z.string().optional(),
        placeholder: z.string().optional(),
        width: z.number().optional(),
      }),
      description: 'A text input field for entering a string value',
    },
    /** Secret text input that hides its value once configured. */
    SecretInput: {
      props: z.object({
        label: z.string(),
        value: z.string().optional(),
        placeholder: z.string().optional(),
        isConfigured: z.boolean().optional(),
        width: z.number().optional(),
      }),
      description: 'A secret input field that masks its value',
    },
    /** Multi-line secret text area (e.g. PEM private keys). */
    SecretTextArea: {
      props: z.object({
        label: z.string(),
        placeholder: z.string().optional(),
        isConfigured: z.boolean().optional(),
        cols: z.number().optional(),
        rows: z.number().optional(),
      }),
      description: 'A secret multi-line text area for entering sensitive multi-line content',
    },
    /** Horizontal divider. */
    Divider: {
      props: z.object({}),
      description: 'A horizontal visual divider',
    },
    /** Data-source description header. */
    Description: {
      props: z.object({
        dataSourceName: z.string(),
        docsLink: z.string(),
      }),
      description: 'Data source description and documentation link header',
    },
  },
  actions: {
    /** Fired when a non-secure jsonData field changes. */
    updateJsonData: {
      params: z.object({ field: z.string(), value: z.string() }),
      description: 'Update a field in jsonData',
    },
    /** Fired when a secure field value changes. */
    updateSecureJsonData: {
      params: z.object({ field: z.string(), value: z.string() }),
      description: 'Update a field in secureJsonData',
    },
    /** Reset a previously-configured secure field. */
    resetSecureField: {
      params: z.object({ field: z.string() }),
      description: 'Reset a configured secure field so the user can re-enter it',
    },
  },
});
