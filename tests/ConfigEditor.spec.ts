import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';

const type = 'grafana-github-datasource';

test('ConfigEditor smoke test', async ({ createDataSourceConfigPage, page, selectors }) => {
  const configPage = await createDataSourceConfigPage({ type });
  page.route(selectors.apis.DataSource.datasourceByUID(configPage.datasource.uid), (route, request) => {
    const data = request.postDataJSON();
    expect(data.jsonData?.githubUrl).toBe('https://github.mycompany.com');
    expect(data.secureJsonData?.accessToken).toBe('my-access-token');
    route.fulfill({ status: 200 });
  });
  await configPage.getByGrafanaSelector(components.ConfigEditor.AccessToken).fill('my-access-token');
  await page.getByRole('radio', { name: 'Enterprise' }).check();
  await page.getByPlaceholder('URL of GitHub Enterprise').fill('https://github.mycompany.com');
  await configPage.saveAndTest();
});
