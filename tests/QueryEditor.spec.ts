import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';
import { githubResponse } from './mocks/github-response';

const type = 'grafana-github-datasource';
let datasourceName = '';

test.beforeAll(async ({ createDataSource }) => {
  const datasource = await createDataSource({ type });
  datasourceName = datasource.name;
});

test('QueryEditor smoke test', async ({ panelEditPage, page }) => {
  await panelEditPage.mockQueryDataResponse(githubResponse);
  await panelEditPage.setVisualization('Table');
  await panelEditPage.datasource.set(datasourceName);
  await panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel).click();
  const select = page.getByLabel('Select options menu');
  await select.locator(page.getByText('Releases')).click();
  await panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
  await panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input).fill('grafana-github-datasource');

  await panelEditPage.refreshPanel();
  try {
    // Newer versions of table view uses gridcell instead of cell
    await expect(page.getByRole('gridcell', { name: 'grafana-github-datasource v1.5.7' })).toBeVisible();
  } catch (error) {
    await expect(page.getByRole('cell', { name: 'grafana-github-datasource v1.5.7' })).toBeVisible();
  }
});
