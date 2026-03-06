import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';

test.describe('Query editor data queries', () => {
  test('Pull_Requests query should return timeseries data', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasource.yaml' });
    await panelEditPage.datasource.set(ds.name);
    await panelEditPage.setVisualization('Time series');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel).click();
    await page.getByLabel('Select options menu').locator(page.getByText('Pull Requests')).click();
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input).fill('plugin-tools');
    await expect(panelEditPage.refreshPanel()).toBeOK();
    await expect(panelEditPage.panel.locator).toBeVisible();
    await expect(panelEditPage.panel.getErrorIcon()).not.toBeVisible();
  });

  test('Pull_Requests table query should return pull request fields', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasource.yaml' });
    await panelEditPage.datasource.set(ds.name);
    await panelEditPage.setVisualization('Table');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel).click();
    await page.getByLabel('Select options menu').locator(page.getByText('Pull Requests')).click();
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input).fill('plugin-tools');
    await expect(panelEditPage.refreshPanel()).toBeOK();
    await expect(panelEditPage.panel.fieldNames).toContainText(['number', 'title', 'url']);
  });

  // TODO: unskip when @grafana/plugin-e2e supports Grafana 12.2.0 viz picker
  test.skip('Pull_Requests table query with search filter should return results', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasource.yaml' });
    await panelEditPage.datasource.set(ds.name);
    await panelEditPage.setVisualization('Table');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel).click();
    await page.getByLabel('Select options menu').locator(page.getByText('Pull Requests')).click();
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input).fill('plugin-tools');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Query.input).fill('is:merged');
    await expect(panelEditPage.refreshPanel()).toBeOK();
    await expect(panelEditPage.panel.fieldNames).toContainText(['number', 'title', 'url']);
  });
});
