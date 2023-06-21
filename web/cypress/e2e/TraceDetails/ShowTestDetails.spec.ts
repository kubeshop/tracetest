import {PokeshopDemo} from '../constants/Test';
import {getResultId, getTestId} from '../utils/Common';

describe('Show test details', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should show the test details for any trace', () => {
    cy.location('pathname').then(pathname => {
      const testId = getTestId(pathname);
      cy.createMultipleTestRuns(testId, 5);

      cy.get(`[data-cy=collapse-test-${testId}]`).click();
      cy.get('[data-cy=test-details-link]', {timeout: 10000}).first().click();

      cy.location('pathname').should('match', /\/test\/.*/i);
      cy.get('[data-cy=test-details-name]').should('have.text', `${PokeshopDemo[0].name} (v1)`);
      cy.get('[data-cy=run-card-list]').should('be.visible');
      cy.get('[data-cy^=run-card-]').should('have.length.above', 0);

      cy.get(`[data-cy=test-details-run-test-button]`).click();
      cy.matchTestRunPageUrl();
    });
  });

  it('should display the test definition yaml', () => {
    cy.location('pathname').then(pathname => {
      const testId = getTestId(pathname);
      cy.visit(`/test/${testId}`);

      cy.get('[data-cy^=result-actions-button]').first().click();
      cy.get('[data-cy=automate-test-button]').click();

      cy.get('[data-cy=code-block]').should('be.visible');
      cy.contains('div', 'Test Definition').should('be.visible');
      cy.contains('div', 'Running Techniques').should('be.visible');
    });
  });

  it('should run a new test', () => {
    cy.location('pathname').then(pathname => {
      const testId = getTestId(pathname);
      const testRunResultId = getResultId(pathname);
      cy.visit(`/test/${testId}`);
      cy.wait('@testObject');
      cy.get(`[data-cy=test-details-run-test-button]`).click();
      cy.matchTestRunPageUrl();
      cy.get('[data-cy=test-header-back-button]').click();
      cy.get(`[data-cy=run-card-${testRunResultId}]`, {timeout: 10000}).should('be.visible');
      cy.visit(`/test/${testId}/run/${testRunResultId}`);
      cy.matchTestRunPageUrl();
    });
  });
  it('should display the jUnit report', () => {
    cy.location('pathname').then(pathname => {
      const testId = getTestId(pathname);
      cy.visit(`/test/${testId}`);
      cy.get('[data-cy^=result-actions-button]').last().click();
      cy.get('[data-cy=view-junit-button]').click();
      cy.get('[data-cy=file-viewer-code-container]').should('be.visible');
      cy.get('[data-cy=file-viewer-close]').click();
      cy.get('[data-cy=file-viewer-code-container]').should('not.be.visible');
      cy.get(`[data-cy=test-details-run-test-button]`).click();
      cy.matchTestRunPageUrl();
    });
  });
});
