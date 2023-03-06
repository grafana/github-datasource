import { e2e } from '@grafana/e2e';

export const openDashboardSettings = (sectionName = 'Variables') => {
  e2e.components.PageToolbar.item('Dashboard settings').click();

  cy.contains(sectionName).should('be.visible').click();
};

export const selectDropdown = (container: Cypress.Chainable<JQuery<HTMLElement>>, text: string, wait = 0) => {
  container.within(() => e2e().get('[class$="-input-suffix"]').click());
  e2e.components.Select.input().should('be.visible').contains(text).click();
  if (wait > 0) {
    e2e().wait(wait);
  }
};
