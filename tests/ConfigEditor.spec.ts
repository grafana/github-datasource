import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';
import semver from 'semver';

const type = 'grafana-github-datasource';

test('ConfigEditor smoke test', async ({ createDataSourceConfigPage, page, selectors, grafanaVersion }) => {
  const configPage = await createDataSourceConfigPage({ type });
  page.route(selectors.apis.DataSource.datasourceByUID(configPage.datasource.uid), (route, request) => {
    const data = request.postDataJSON();
    expect(data.jsonData?.githubUrl).toBe('https://github.mycompany.com');
    expect(data.secureJsonData?.accessToken).toBe('my-access-token');
    route.fulfill({ status: 200 });
  });
  await configPage.getByGrafanaSelector(components.ConfigEditor.AccessToken).fill('my-access-token');
  // TODO: Move this logic to plugin-e2e
  if (semver.lt(grafanaVersion, '10.2.0')) {
    await page.getByText('Enterprise Server', { exact: true }).click();
  } else {
    await page.getByRole('radio', { name: 'Enterprise Server' }).check();
  }
  await page.getByPlaceholder('http(s)://HOSTNAME/').fill('https://github.mycompany.com');
  // TODO: Move this logic to plugin-e2e
  if (semver.lt(grafanaVersion, '10.0.0')) {
    await page.getByLabel('Data source settings page Save and Test button').click();
  } else {
    await configPage.saveAndTest();
  }
});
