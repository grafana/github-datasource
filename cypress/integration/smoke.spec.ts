import { e2e } from '@grafana/e2e';
import { PartialConfigurePanelConfig } from '@grafana/e2e/flows/configurePanel';
import { selectors } from '../../src/components/selectors';

const e2eSelectors = e2e.getSelectors(selectors.components);

const addGithubDataSource = (accessToken: string) => {
  return e2e.flows.addDataSource({
    checkHealth: true,
    expectedAlertMessage: 'OK',
    form: () => {
      e2eSelectors.ConfigEditor.AccessToken.input()
        .scrollIntoView()
        .type(accessToken);
    },
    type: 'GitHub',
  });
};

const addGithubPanel = (variableName: string) => {
  const fillQueryEditor = () => {
    // Fill in the Query Editor with data needed
    // We are specifically looking at Releases with `grafana` repo to ensure we don't have any
    // sensitive info that will be caught by the screenshot
    e2e.components.QueryEditorRows.rows()
      .should('be.visible')
      .within(() => {
        e2e.components.Select.input()
          .first()
          .should('be.empty')
          .focus()
          .type(`Releases{enter}`);
      });

    e2eSelectors.QueryEditor.Owner.input()
      .should('be.empty')
      .type(`grafana{enter}`);
    e2eSelectors.QueryEditor.Repository.input()
      .should('be.empty')
      .type(`grafana{enter}`);
  };

  e2e.flows
    .addPanel({
      matchScreenshot: true,
      queriesForm: () => {
        fillQueryEditor();

        // Switch to the Table view as we have tabular data from Github (for the screenshot)
        e2e.components.PanelEditor.OptionsPane.open()
          .should('be.visible')
          .click();
        e2e.components.OptionsGroup.toggle('Panel type')
          .should('be.visible')
          .click();
        e2e.components.PluginVisualization.item('Table')
          .scrollIntoView()
          .should('be.visible')
          .click();
      },
    })
    .then(({ config: { panelTitle } }: { config: PartialConfigurePanelConfig }) => {
      // Make sure the template variable works
      e2e.flows.editPanel({
        matchScreenshot: true,
        panelTitle,
        queriesForm: () => {
          e2eSelectors.QueryEditor.Owner.input()
            .clear()
            .type(`$${variableName}{enter}`);
        },
        visitDashboardAtStart: false,
      });
    });
};

e2e.scenario({
  describeName: 'Smoke tests',
  itName: 'Login, create data source, annotation, template variable, dashboard, panel',
  scenario: () => {
    e2e()
      .readProvisions([
        // Paths are relative to <project-root>/provisioning
        'datasources/github.yaml',
      ])
      .then(([provision]) => addGithubDataSource(provision.datasources[0].secureJsonData.accessToken))
      .then(({ config: { name: dataSourceName } }) => {
        const variableName = 'owner';

        e2e.flows.addDashboard({
          annotations: [
            {
              dataSource: dataSourceName,
              dataSourceForm: () => e2eSelectors.AnnotationEditor.container().should('be.visible'),
              name: 'Annotations',
            },
          ],
          timeRange: {
            from: '2020-12-01 00:00:00',
            to: '2020-12-31 23:59:59',
          },
          variables: [
            {
              constantValue: 'grafana',
              label: 'Template Variable',
              name: variableName,
              type: e2e.flows.VARIABLE_TYPE_CONSTANT,
            },
          ],
        });

        addGithubPanel(variableName);
      });
  },
});
