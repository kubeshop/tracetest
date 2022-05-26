import {createMultipleTestRuns, createTest, deleteTest, description, getResultId, name, testId} from '../utils/common';

describe('Show test details', () => {
  before(() => {
    createTest();
  });

  after(() => {
    deleteTest();
  });

  it('should show the test details for any trace', () => {
    createMultipleTestRuns(testId, 5);
    cy.get(`[data-cy=collapse-test-${testId}]`).click();
    cy.get('[data-cy=test-details-link]', {timeout: 10000}).first().click();

    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} - ${description} (v1)`);
    cy.get('[data-cy=result-card-list]').should('be.visible');
    cy.get('[data-cy^=result-card-]').should('have.length.above', 0);
  });

  it('should run a new test', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);
    cy.get(`[data-cy=test-details-run-test-button]`).click();
    cy.location('pathname').should('match', /\/run\/.*/i);

    cy.location().then(({pathname}) => {
      const testRunResultId = getResultId(pathname);

      cy.wait(2000);
      cy.get('[data-cy=test-header-back-button]').click();
      cy.get(`[data-cy=result-card-${testRunResultId}]`, {timeout: 10000}).should('be.visible');
      cy.visit(`http://localhost:3000/test/${testId}/run/${testRunResultId}`);
    });
  });

  // it('should update the test run result status', () => {
  //   cy.visit(`http://localhost:3000/test/${testId}`);
  //   cy.get(`[data-cy=test-details-run-test-button]`).click();
  //   cy.location('pathname').should('match', /\/result\/.*/i);

  //   cy.location().then(({pathname}) => {
  //     const testRunResultId = getResultId(pathname);

  //     cy.get('[data-cy=test-run-result-status]', { timeout: 2000 }).should('have.text', 'Test status:Awaiting trace');
  //     cy.wait(2000);
  //     cy.get('[data-cy=test-run-result-status]').should('have.text', 'Test status:Finished');
  //     cy.get('[data-cy=test-header-back-button]').click();
  //     cy.get(`[data-cy=test-run-result-status-${testRunResultId}]`, { timeout: 10000 }).should('have.text', 'Finished');
  //     cy.visit(`http://localhost:3000/test/${testId}/result/${testRunResultId}`);
  //   });
  // });
});
