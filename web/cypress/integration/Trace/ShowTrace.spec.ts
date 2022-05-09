import { createTest, deleteTest, testId } from "../utils/common";

describe('Show Trace', () => {
  beforeEach(() => {
    createTest();
  });

  afterEach(() => {
    deleteTest();
  });

  it('should show the trace components', () => {
    cy.get(`[data-cy=test-url-${testId}]`).first().click();
    cy.location('href').should('match', /\/test\/.*/i);

    cy.get(`[data-cy^=test-run-result-]`).first().click();
    cy.location('href').should('match', /resultId=.*/i);

    cy.wait(6000);
    cy.get('[data-cy=diagram-dag]').should('be.visible');
    cy.get('[data-cy^=trace-node-]').should('be.visible');
    cy.get('[data-cy=span-details-attributes]').should('be.visible');
    cy.get('[data-cy=empty-assertion-table]').should('be.visible');

    cy.get('#rc-tabs-1-tab-test-results').click();
    cy.get('[data-cy=test-results]').should('be.visible');
  });
});
