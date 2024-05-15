import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';
import { githubResponse } from './mocks/github-response';

const type = 'grafana-github-datasource';

test('QueryEditor', async ({ panelEditPage, page, createDataSource }) => {
  const datasource = await createDataSource({ type });
  await panelEditPage.mockQueryDataResponse(githubResponse);
  await panelEditPage.setVisualization('Table');
  await panelEditPage.datasource.set(datasource.name);
  await panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel).click();
  const select = page.getByLabel('Select options menu');
  await select.locator(page.getByText('Releases')).click();
  await panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
  await panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input).fill('grafana-github-datasource');

  await panelEditPage.refreshPanel();
  await expect(page.getByRole('cell', { name: 'grafana-github-datasource v1.5.7' })).toBeVisible();
});

test('ConfigEditor', async ({ createDataSourceConfigPage, page, selectors }) => {
  const configPage = await createDataSourceConfigPage({ type });
  page.route(selectors.apis.DataSource.datasourceByUID(configPage.datasource.uid), (route, request) => {
    const data = request.postDataJSON();
    expect(data.jsonData?.githubUrl).toBe('https://github.mycompany.com');
    expect(data.secureJsonData?.accessToken).toBe('my-access-token');
    route.fulfill({ status: 200 });
  });
  await configPage.getByGrafanaSelector(components.ConfigEditor.AccessToken).fill('my-access-token');
  await page.getByLabel('Enterprise').click();
  await page.getByPlaceholder('URL of GitHub Enterprise').fill('https://github.mycompany.com');
  await configPage.saveAndTest();
});
