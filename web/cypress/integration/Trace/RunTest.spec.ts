import {createTest, deleteTest, extractTestRunIdFromTracePage} from '../utils/Common';

describe('Run Test', () => {
  it('should show and click the Run Test button when the test has finished', () => {
    (async () => {
      const testId = await createTest();
      cy.visit(`http://localhost:3000/test/${testId}`);
      cy.get('[data-cy^=result-card]', {timeout: 10000}).first().click();
      cy.location('href').should('match', /\/test\/.*/i);

      cy.get(`[data-cy^=test-run-result-]`).first().click();
      cy.location('href').should('match', /\/run\/.*/i);

      cy.get('[data-cy=run-test-button]', {timeout: 20000}).should('be.visible');
      cy.get(`[data-cy^=run-test-button]`).first().click();

      const testRunResultId = await extractTestRunIdFromTracePage();
      cy.location('pathname').should('eq', `/test/${testId}/run/${testRunResultId}`);
      deleteTest(testId);
    })();
  });
});
