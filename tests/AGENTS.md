# E2E testing a Grafana plugin

This plugin uses `@grafana/plugin-e2e` and Playwright for end-to-end testing.

- Always import `test` and `expect` from `@grafana/plugin-e2e`, not from `@playwright/test`.
- Always use `@grafana/plugin-e2e` fixtures and page models instead of raw Playwright navigation. They handle Grafana version differences automatically.
- Place test files in `tests/` as `*.spec.ts`.
- Each test must be independent and assume fresh state.
- If tests fail against newer Grafana versions, update `@grafana/plugin-e2e` first. It evolves alongside Grafana core to handle selector and API changes.

## Selecting elements

### Grafana selectors

- Use Grafana e2e-selectors whenever possible. Always get them from the `selectors` fixture provided by `@grafana/plugin-e2e` - never import from `@grafana/e2e-selectors` directly. The fixture resolves the correct selectors for the Grafana version under test.
- Always use the `getByGrafanaSelector` method (exposed by all plugin-e2e page models) to resolve selectors to Playwright locators. It handles the `aria-label` vs `data-testid` difference across Grafana versions automatically.
  ```typescript
  panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container).click();
  ```

### Scoping locators

Scope locators to the narrowest context possible.

```typescript
// bad - matches any "URL" text on the page
page.getByText('URL').click();
// good - scoped to the plugin's wrapper
page.getByTestId('plugin-url-wrapper').getByText('URL').click();
```

### Form elements

The `InlineField` and `Field` components can be used interchangeably in the examples below.

**Input** - use `getByRole('textbox', { name: '<label>' })` where the name matches the wrapping `InlineField` label.

```tsx
// component
<InlineField label="Auth key">
  <Input value={value} onChange={handleOnChange} id="config-auth-key" />
</InlineField>
```

```typescript
// test
await page.getByRole('textbox', { name: 'Auth key' }).fill('..');
```

**Select** - use `getByRole('combobox', { name: '<label>' })` to open the dropdown. Assert options using `selectors.components.Select.option` via `getByGrafanaSelector`. Note: the `Select` component requires `inputId` (not `id`) for label association.

```tsx
// component
<InlineField label="Auth type">
  <Select inputId="config-auth-type" value={value} options={options} onChange={handleOnChange} />
</InlineField>
```

```typescript
// test
await page.getByRole('combobox', { name: 'Auth type' }).click();
const option = selectors.components.Select.option;
await expect(configPage.getByGrafanaSelector(option)).toHaveText(['val1', 'val2']);
```

**Checkbox** - use `getByRole('checkbox', { name: '<label>' })`. The underlying input is not directly clickable so you must pass `{ force: true }`.

```tsx
// component
<InlineField label="TLS Enabled">
  <Checkbox id="config-tls-enabled" value={value} onChange={handleOnChange} />
</InlineField>
```

```typescript
// test
await page.getByRole('checkbox', { name: 'TLS Enabled' }).uncheck({ force: true });
await expect(page.getByRole('checkbox', { name: 'TLS Enabled' })).not.toBeChecked();
```

**InlineSwitch** - use `getByLabel('<label>')`. Like Checkbox, requires `{ force: true }`. The `InlineSwitch` `label` prop must match the wrapping `InlineField` label.

```tsx
// component
<InlineField label="TLS Enabled">
  <InlineSwitch label="TLS Enabled" value={value} onChange={handleOnChange} />
</InlineField>
```

```typescript
// test
await page.getByLabel('TLS Enabled').uncheck({ force: true });
await expect(page.getByLabel('TLS Enabled')).not.toBeChecked();
```

## Using the plugin-e2e API

`@grafana/plugin-e2e` exposes page models and fixtures that encapsulate common UI operations and handle Grafana version differences.

To discover all available fixtures, options, models and matchers, read the exports in `node_modules/@grafana/plugin-e2e/src/index.ts`.

### Fixtures

Fixtures follow a naming convention that indicates how the resource is obtained:

- **camelCase** (e.g. `panelEditPage`) - creates a new, empty resource. Use when testing from a blank state.
- **`goto` prefix** (e.g. `gotoPanelEditPage`) - navigates to an existing resource. Use with provisioned dashboards and datasources.
- **`readProvisioned` prefix** (e.g. `readProvisionedDataSource`) - reads a provisioning file from disk. Use to avoid hardcoding UIDs and names.

### Options

Default options like `featureToggles`, `user`, `userPreferences` and `provisioningRootDir` can be overridden per project in `playwright.config.ts` via `use`. Per-test overrides work the same way using `test.use()`.

{{#if_eq pluginType "datasource"}}

### Custom matchers

`@grafana/plugin-e2e` extends Playwright's `expect` with Grafana-specific assertions:

- `toBeOK()` - asserts a response status is 200-299. Use with `saveAndTest()`, `refreshPanel()`, `runQuery()` and `evaluate()`.
- `toHaveAlert(severity, options?)` - asserts an alert box is visible. Severity is `'success'`, `'warning'`, `'error'` or `'info'`. Supports `hasText` filtering.
- `toDisplayPreviews(values)` - asserts variable query preview matches expected values. Use with `variableEditPage`.

```typescript
await expect(configPage.saveAndTest()).toBeOK();
await expect(configPage).toHaveAlert('error', { hasText: 'API key is missing' });
await expect(variableEditPage).toDisplayPreviews(['value1', 'value2']);
```

{{/if_eq}}

{{#if_eq pluginType "panel"}}

## Panel options

Prefer provisioning dashboards with panels pre-configured in different states and navigating to them with `gotoPanelEditPage`. This avoids coupling tests to the Grafana panel edit UI and makes them more stable across versions.

When you need to interact with panel options directly, use the option group helpers on `panelEditPage`. For Grafana-provided groups use `getPanelOptions()`, `getStandardOptions()`, `getValueMappingOptions()`, `getDataLinksOptions()` or `getThresholdsOptions()`. For custom groups use `getCustomOptions('Group Name')`.

Each option group returns an object with typed accessors: `getSwitch(label)`, `getSelect(label)`, `getMultiSelect(label)`, `getRadioGroup(label)`, `getTextInput(label)`, `getNumberInput(label)`, `getSliderInput(label)`, `getColorPicker(label)` and `getUnitPicker(label)`.

```typescript
test('should update unit when standard option changes', async ({ panelEditPage }) => {
  const standardOptions = panelEditPage.getStandardOptions();
  await standardOptions.getUnitPicker('Unit').selectOption('Misc > Pixels');
  await expect(panelEditPage.panel.locator).toContainText('px');
});

test('should change timezone when custom option is selected', async ({ panelEditPage, page }) => {
  const tzOptions = panelEditPage.getCustomOptions('Timezone');
  await tzOptions.getSelect('Timezone').selectOption('Europe/Stockholm');
  await expect(page.getByTestId('time-zone')).toContainText('Europe/Stockholm');
});
```

{{/if_eq}}

## Resources

Tests often need pre-configured datasources, dashboards or alert rules. Use [Grafana provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/) to set these up - place YAML/JSON files in the `provisioning/` folder and they will be loaded when the Grafana test server starts.

- Never hardcode UIDs or names. Use `readProvisionedDataSource`, `readProvisionedDashboard` and `readProvisionedAlertRule` fixtures to read values from provisioning files.
- Each test should be independent. Provision the resources you need rather than relying on state from previous tests.
- Use [environment variable interpolation](https://grafana.com/docs/grafana/latest/administration/provisioning/#using-environment-variables) for secrets - never commit credentials to the repository.
- If CI requires provisioning, make sure `provisioning/` is not in `.gitignore`.

## Running tests

Tests should be run against multiple Grafana versions. Check `grafanaDependency` in `src/plugin.json` for the minimum supported version.

**Min supported version**:
Align with `grafanaDependency` in `plugin.json`.

```bash
# terminal 1
GRAFANA_VERSION=11.3.0 yarn run server

# terminal 2
yarn run e2e
```

**Latest dev image** (forwards compatibility):

```bash
# terminal 1
GRAFANA_IMAGE=grafana-dev GRAFANA_VERSION=12.4.0-211043112277 yarn run server

# terminal 2
yarn run e2e
```
