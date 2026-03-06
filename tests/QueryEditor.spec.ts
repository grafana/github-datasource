import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';

// use the same absolute time range as the provisioned dashboard
// so WireMock stubs match in replay mode
const timeRange = { from: '2026-02-05 11:00:00', to: '2026-02-19 10:59:58', zone: 'Coordinated Universal Time' };

test.describe('Query editor data queries', () => {
  test('Pull_Requests query should return timeseries data', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasource.yaml' });
    await panelEditPage.datasource.set(ds.name);
    await panelEditPage.timeRange.set(timeRange);
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
    await panelEditPage.timeRange.set(timeRange);
    await panelEditPage.setVisualization('Table');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel).click();
    await page.getByLabel('Select options menu').locator(page.getByText('Pull Requests')).click();
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
    await panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input).fill('plugin-tools');
    await expect(panelEditPage.refreshPanel()).toBeOK();
    await expect(panelEditPage.panel.fieldNames).toContainText(['number', 'title', 'url']);
  });

  test('Pull_Requests table query with search filter should return results', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasource.yaml' });
    await panelEditPage.datasource.set(ds.name);
    await panelEditPage.timeRange.set(timeRange);
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
