import {createMultipleTestRuns, createTest, deleteTest, getResultId, name} from '../utils/Common';

describe('Show test details', () => {
  let testId;
  before(() => {
    testId = createTest();
  });

  after(() => {
    deleteTest();
  });

  it('should show the test details for any trace', () => {
    createMultipleTestRuns(testId, 5);

    cy.get(`[data-cy=collapse-test-${testId}]`).click();
    cy.get('[data-cy=test-details-link]', {timeout: 20000}).first().click();

    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
    cy.get('[data-cy=result-card-list]').should('be.visible');
    cy.get('[data-cy^=result-card-]').should('have.length.above', 0);
  });

  it('should display the jUnit report', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);

    cy.get('[data-cy^=result-actions-button]').last().click();
    cy.wait(5000);
    cy.get('[data-cy=view-junit-button]').click();

    cy.get('[data-cy=file-viewer-code-container]', {timeout: 20000}).should('be.visible');
    cy.get('[data-cy=file-viewer-close]').click();

    cy.get('[data-cy=file-viewer-code-container]').should('not.be.visible');
  });

  it('should display the test definition yaml', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);

    cy.get('[data-cy^=result-actions-button]').first().click();
    cy.get('[data-cy=view-test-definition-button]').click();

    cy.get('[data-cy=file-viewer-code-container]').should('be.visible');
    cy.get('[data-cy=file-viewer-close]').click();

    cy.get('[data-cy=file-viewer-code-container]').should('not.be.visible');
  });

  it('should run a new test', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);
    cy.get(`[data-cy=test-details-run-test-button]`).click();
    cy.location('pathname').should('match', /\/run\/.*/i);

    cy.location().then(({pathname}) => {
      const testRunResultId = getResultId(pathname);

      cy.wait(2000);
      cy.get('[data-cy=test-header-back-button]').click();
      cy.get(`[data-cy=result-card-${testRunResultId}]`, {timeout: 20000}).should('be.visible');
      cy.visit(`http://localhost:3000/test/${testId}/run/${testRunResultId}`);
    });
  });
});
