import { test, expect } from '@grafana/plugin-e2e';

const dashboardFileName = 'example.json';

test.describe('Dashboard panels', () => {
  test('pull requests timeseries panel should render data', async ({
    readProvisionedDashboard,
    gotoDashboardPage,
  }) => {
    const dashboard = await readProvisionedDashboard({ fileName: dashboardFileName });
    const dashboardPage = await gotoDashboardPage(dashboard);
    const panel = dashboardPage.getPanelByTitle('Pull requests');
    await expect(panel.locator).toBeVisible();
    await expect(panel.getErrorIcon()).not.toBeVisible();
  });

  test('table panel should display pull request data', async ({
    readProvisionedDashboard,
    gotoDashboardPage,
  }) => {
    const dashboard = await readProvisionedDashboard({ fileName: dashboardFileName });
    const dashboardPage = await gotoDashboardPage(dashboard);
    const panel = dashboardPage.getPanelById('3');
    await panel.locator.scrollIntoViewIfNeeded();
    await expect(panel.fieldNames).toContainText(['number', 'title', 'url']);
    await expect(panel.getErrorIcon()).not.toBeVisible();
  });

  test('table panel with variable filter should display data', async ({
    readProvisionedDashboard,
    gotoDashboardPage,
  }) => {
    const dashboard = await readProvisionedDashboard({ fileName: dashboardFileName });
    const dashboardPage = await gotoDashboardPage(dashboard);
    const panel = dashboardPage.getPanelById('2');
    await panel.locator.scrollIntoViewIfNeeded();
    await expect(panel.fieldNames).toContainText(['number', 'title', 'url']);
    await expect(panel.getErrorIcon()).not.toBeVisible();
  });
});
