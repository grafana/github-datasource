import { e2e } from '@grafana/e2e';
import { selectors } from '../../src/components/selectors';
import { openDashboardSettings } from './utils';

const e2eSelectors = e2e.getSelectors(selectors.components);

const addGithubDataSource = (accessToken: string) => {
  return e2e.flows.addDataSource({
    expectedAlertMessage: 'Data source is working',
    form: () => e2eSelectors.ConfigEditor.AccessToken.input().scrollIntoView().type(accessToken),
    timeout: 10000,
    type: 'GitHub',
  });
};

describe('Smoke test', () => {
  it('should render default variable editor', () => {
    e2e.flows.login();
    e2e()
      .readProvisions(['datasources/github.yaml'])
      .then(([provision]) => addGithubDataSource(provision.datasources[0].secureJsonData.accessToken))
      .then(({ config: { name: dataSourceName } }) => {
        e2e.flows
          .addDashboard({
            timeRange: {
              from: '2020-12-01 00:00:00',
              to: '2020-12-31 23:59:59',
            },
          })
          .then(() => {
            const variableName = 'owner';
            openDashboardSettings('Variables');
            e2e.pages.Dashboard.Settings.Variables.List.addVariableCTAV2().click();
            e2e.pages.Dashboard.Settings.Variables.Edit.General.generalNameInputV2().clear().type(variableName);
            cy.wait(6 * 1000); // When clearing the variable name, the validation popup comes and hides the datasource picker. so wait sometime till the popup closes.
            e2e.pages.Dashboard.Settings.Variables.Edit.General.generalTypeSelectV2().within(() => {
              e2e().get('input').type('Constant{enter}');
            });
            e2e.pages.Dashboard.Settings.Variables.Edit.ConstantVariable.constantOptionsQueryInputV2()
              .type('grafana')
              .blur();
            e2e.pages.Dashboard.Settings.Variables.Edit.General.previewOfValuesOption()
              .eq(0)
              .should('have.text', 'grafana');
            e2e.pages.Dashboard.Settings.Variables.Edit.General.submitButton().click();
            e2e.components.BackButton.backArrow().click({ force: true });
            e2e.components.RefreshPicker.runButtonV2().click();
            // e2e.flows.addPanel({
            //   queriesForm: () => {},
            // });
          });
      });
  });
});
