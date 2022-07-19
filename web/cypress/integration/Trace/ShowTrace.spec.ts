import {createTest, deleteTest} from '../utils/Common';

describe('Show Trace', () => {
  let testId;
  beforeEach(() => {
    testId = createTest();
  });

  afterEach(() => {
    deleteTest();
  });

  it('should show the trace components', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);
    cy.get('[data-cy^=result-card]', {timeout: 20000}).first().click();

    cy.location('href').should('match', /\/test\/.*/i);

    cy.get(`[data-cy^=test-run-result-]`).first().click();
    cy.location('href').should('match', /\/run\/.*/i);

    cy.get('[data-cy^=trace-node-]', {timeout: 30000}).should('be.visible');
    cy.get('[data-cy=span-details-attributes]').should('be.visible');
    cy.get('[data-cy=empty-assertion-card-list]').should('exist');

    cy.get('[data-cy=assertion-card-list]').should('exist');
  });
});
