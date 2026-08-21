import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';
import { githubVariableResponse } from './mocks/github-response';

const type = 'grafana-github-datasource';
let datasourceName = '';

test.beforeAll(async ({ createDataSource }) => {
  const datasource = await createDataSource({ type });
  datasourceName = datasource.name;
});

test('Variable query editor renders values from a GitHub query', async ({ variableEditPage, page }) => {
  await variableEditPage.mockQueryDataResponse(githubVariableResponse);
  await variableEditPage.setVariableType('Query');
  await variableEditPage.datasource.set(datasourceName);

  const queryTypeContainer = variableEditPage.getByGrafanaSelector(
    components.QueryEditor.QueryType.container.ariaLabel
  );
  await queryTypeContainer.getByRole('combobox').click();
  await page.getByRole('listbox').getByRole('option', { name: 'Releases' }).click();
  await variableEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input).fill('grafana');
  const repositoryInput = variableEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input);
  await repositoryInput.fill('github-datasource');
  await repositoryInput.blur();

  const queryDataRequest = variableEditPage.waitForQueryDataRequest();
  await variableEditPage.runQuery();
  const request = await queryDataRequest;
  expect(request.postDataJSON().queries[0]).toMatchObject({
    refId: 'metricFindQuery',
    queryType: 'Releases',
    owner: 'grafana',
    repository: 'github-datasource',
  });
  await expect(variableEditPage).toDisplayPreviews(['grafana-github-datasource v1.5.7']);
});
