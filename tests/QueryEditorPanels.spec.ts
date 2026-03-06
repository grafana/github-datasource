import { test, expect } from '@grafana/plugin-e2e';
import { components } from '../src/components/selectors';

const dashboardFileName = 'example.json';

test.describe('Query editor in provisioned panels', () => {
  test('panel 1 - should load Pull_Requests query type with owner and repository', async ({
    readProvisionedDashboard,
    gotoPanelEditPage,
  }) => {
    const dashboard = await readProvisionedDashboard({ fileName: dashboardFileName });
    const panelEditPage = await gotoPanelEditPage({ dashboard, id: '1' });
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel)).toContainText(
      'Pull Requests'
    );
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input)).toHaveValue('grafana');
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input)).toHaveValue(
      'plugin-tools'
    );
  });

  test('panel 2 - should load Pull_Requests query with title variable filter', async ({
    readProvisionedDashboard,
    gotoPanelEditPage,
  }) => {
    const dashboard = await readProvisionedDashboard({ fileName: dashboardFileName });
    const panelEditPage = await gotoPanelEditPage({ dashboard, id: '2' });
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel)).toContainText(
      'Pull Requests'
    );
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input)).toHaveValue('grafana');
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input)).toHaveValue(
      'plugin-tools'
    );
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.Query.input)).toHaveValue('title:"$title"');
  });

  test('panel 3 - should load Pull_Requests query with owner and repository', async ({
    readProvisionedDashboard,
    gotoPanelEditPage,
  }) => {
    const dashboard = await readProvisionedDashboard({ fileName: dashboardFileName });
    const panelEditPage = await gotoPanelEditPage({ dashboard, id: '3' });
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.QueryType.container.ariaLabel)).toContainText(
      'Pull Requests'
    );
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.Owner.input)).toHaveValue('grafana');
    await expect(panelEditPage.getByGrafanaSelector(components.QueryEditor.Repository.input)).toHaveValue(
      'plugin-tools'
    );
  });
});
