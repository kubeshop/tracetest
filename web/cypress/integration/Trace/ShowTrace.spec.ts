import {createTest, deleteTest, testId} from '../utils/common';

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
    cy.location('href').should('match', /\/result\/.*/i);

    cy.get('[data-cy=diagram-dag]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy^=trace-node-]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=span-details-attributes]').should('be.visible');
    cy.get('[data-cy=empty-assertion-table]').should('be.visible');

    cy.get('[data-cy=test-results]').should('exist');
  });
});
