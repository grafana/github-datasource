import React from 'react';
import { defineRegistry, useBoundProp } from '@json-render/react';
import { ConfigSection, DataSourceDescription } from '@grafana/plugin-ui';
import {
  Field,
  Input,
  Label,
  RadioButtonGroup,
  SecretInput as GrafanaSecretInput,
  SecretTextArea as GrafanaSecretTextArea,
} from '@grafana/ui';
import { Divider } from 'components/Divider';
import { configEditorCatalog } from './catalog';

/**
 * Create a json-render registry that maps catalog component names to
 * concrete React implementations backed by Grafana UI primitives.
 *
 * The registry also wires up catalog actions so that field mutations
 * can be forwarded to the Grafana DataSourcePluginOptionsEditor callback.
 *
 * @param onJsonDataChange  called when a non-secure jsonData field changes
 * @param onSecureJsonDataChange  called when a secure field value changes
 * @param onResetSecureField  called when a configured secret needs resetting
 */
export function createConfigEditorRegistry(
  onJsonDataChange: (field: string, value: string) => void,
  onSecureJsonDataChange: (field: string, value: string) => void,
  onResetSecureField: (field: string) => void
) {
  return defineRegistry(configEditorCatalog, {
    components: {
      Stack: ({ children }) => <div>{children}</div>,

      Section: ({ props, children }) => (
        <ConfigSection title={props.title} isCollapsible={props.isCollapsible}>
          {children}
        </ConfigSection>
      ),

      RadioGroup: ({ props, bindings }) => {
        const [value, setValue] = useBoundProp(props.value, bindings?.value);
        return (
          <>
            {props.label && <Label>{props.label}</Label>}
            <RadioButtonGroup
              options={props.options}
              value={value ?? ''}
              onChange={(v) => {
                setValue(v);
                // Also propagate via action so external listeners can react
                const fieldPath = bindings?.value;
                if (fieldPath) {
                  onJsonDataChange(fieldPath.replace(/^\//, ''), v);
                }
              }}
            />
          </>
        );
      },

      TextInput: ({ props, bindings }) => {
        const [value, setValue] = useBoundProp(props.value, bindings?.value);
        return (
          <Field label={props.label}>
            <Input
              placeholder={props.placeholder}
              value={value ?? ''}
              width={props.width ?? 40}
              onChange={(e) => {
                const v = e.currentTarget.value;
                setValue(v);
                const fieldPath = bindings?.value;
                if (fieldPath) {
                  onJsonDataChange(fieldPath.replace(/^\//, ''), v);
                }
              }}
            />
          </Field>
        );
      },

      SecretInput: ({ props }) => {
        const isConfigured = props.isConfigured ?? false;
        return (
          <Field label={props.label}>
            <GrafanaSecretInput
              placeholder={props.placeholder}
              value={props.value ?? ''}
              isConfigured={isConfigured}
              width={props.width ?? 40}
              onChange={(e) => {
                onSecureJsonDataChange(
                  props.label.replace(/\s+/g, '').replace(/^(.)/, (c) => c.toLowerCase()),
                  e.currentTarget.value
                );
              }}
              onReset={() => {
                onResetSecureField(
                  props.label.replace(/\s+/g, '').replace(/^(.)/, (c) => c.toLowerCase())
                );
              }}
            />
          </Field>
        );
      },

      SecretTextArea: ({ props }) => {
        const isConfigured = props.isConfigured ?? false;
        return (
          <Field label={props.label}>
            <GrafanaSecretTextArea
              placeholder={props.placeholder}
              isConfigured={isConfigured}
              cols={props.cols ?? 55}
              rows={props.rows ?? 7}
              onChange={(e) => {
                onSecureJsonDataChange(
                  props.label.replace(/\s+/g, '').replace(/^(.)/, (c) => c.toLowerCase()),
                  e.currentTarget.value
                );
              }}
              onReset={() => {
                onResetSecureField(
                  props.label.replace(/\s+/g, '').replace(/^(.)/, (c) => c.toLowerCase())
                );
              }}
            />
          </Field>
        );
      },

      Divider: () => <Divider />,

      Description: ({ props }) => (
        <DataSourceDescription
          dataSourceName={props.dataSourceName}
          docsLink={props.docsLink}
          hasRequiredFields={false}
        />
      ),
    },
    actions: {
      updateJsonData: async (params, _setState, _state) => {
        if (params) {
          onJsonDataChange(params.field, params.value);
        }
      },
      updateSecureJsonData: async (params) => {
        if (params) {
          onSecureJsonDataChange(params.field, params.value);
        }
      },
      resetSecureField: async (params) => {
        if (params) {
          onResetSecureField(params.field);
        }
      },
    },
  });
}
