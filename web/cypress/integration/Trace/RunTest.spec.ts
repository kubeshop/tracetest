import {createTest, deleteTest, getResultId, testId} from '../utils/Common';

describe('Run Test', () => {
  beforeEach(() => {
    createTest();
  });

  afterEach(() => {
    deleteTest();
  });

  it('should show and click the Run Test button when the test has finished', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);
    cy.get('[data-cy^=result-card]', {timeout: 20000}).first().click();
    cy.location('href').should('match', /\/test\/.*/i);

    cy.get(`[data-cy^=test-run-result-]`).first().click();
    cy.location('href').should('match', /\/run\/.*/i);

    cy.get('[data-cy=run-test-button]', {timeout: 20000}).should('be.visible');
    cy.get(`[data-cy^=run-test-button]`).first().click();

    cy.wait(2000);
    cy.location().then(({pathname}) => {
      const testRunResultId = getResultId(pathname);
      cy.location('pathname').should('eq', `/test/${testId}/run/${testRunResultId}`);
    });
  });
});
