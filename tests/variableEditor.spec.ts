import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';

test.describe('Variable queries', () => {
  test('provisioned variable should load query editor with correct fields', async ({
    readProvisionedDashboard,
    gotoVariableEditPage,
    selectors,
  }) => {
    const dashboard = await readProvisionedDashboard({ fileName: 'example.json' });
    const variableEditPage = await gotoVariableEditPage({ dashboard, id: '0' });
    await expect(
      variableEditPage.getByGrafanaSelector(selectors.pages.Dashboard.Settings.Variables.Edit.General.generalNameInputV2)
    ).toHaveValue('title');
    await expect(variableEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel)).toContainText(
      'Pull Requests'
    );
    await expect(variableEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input)).toHaveValue('grafana');
    await expect(variableEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input)).toHaveValue(
      'plugin-tools'
    );
  });

  test('new Pull_Requests variable query should run without error', async ({
    variableEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasource.yaml' });
    await variableEditPage.setVariableType('Query');
    await variableEditPage.datasource.set(ds.name);
    await variableEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel).click();
    await page.getByLabel('Select options menu').locator(page.getByText('Pull Requests')).click();
    await variableEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
    await variableEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input).fill('plugin-tools');
    const queryDataRequest = variableEditPage.waitForQueryDataRequest();
    await variableEditPage.runQuery();
    await queryDataRequest;
  });
});
