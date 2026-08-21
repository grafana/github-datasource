import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';
import { githubBranchesResponse, githubResponse } from './mocks/github-response';

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
  const queryTypeContainer = panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel);
  await queryTypeContainer.getByRole('combobox').click();
  await page.getByRole('listbox').getByRole('option', { name: 'Releases' }).click();
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

test('Branches query sends its filter and renders the response', async ({ panelEditPage, page }) => {
  await panelEditPage.mockQueryDataResponse(githubBranchesResponse);
  await panelEditPage.setVisualization('Table');
  await panelEditPage.datasource.set(datasourceName);

  const queryTypeContainer = panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel);
  await queryTypeContainer.getByRole('combobox').click();
  await page.getByRole('listbox').getByRole('option', { name: 'Branches' }).click();
  await panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
  await panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input).fill('github-datasource');
  await page.getByRole('textbox', { name: 'Branch filter' }).fill('release/');
  await page.getByRole('textbox', { name: 'Branch filter' }).blur();

  const queryDataRequest = panelEditPage.waitForQueryDataRequest();
  await expect(panelEditPage.refreshPanel()).toBeOK();
  const request = await queryDataRequest;
  expect(request.postDataJSON().queries[0]).toMatchObject({
    queryType: 'Branches',
    owner: 'grafana',
    repository: 'github-datasource',
    options: { query: 'release/' },
  });
  await expect(
    page.getByRole('gridcell', { name: 'release/2.9' }).or(page.getByRole('cell', { name: 'release/2.9' }))
  ).toBeVisible();
});
