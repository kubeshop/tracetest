import {createTest, deleteTest, testId} from '../utils/Common';

describe('Edit Test', () => {
  before(() => {
    createTest();
  });

  after(() => {
    deleteTest();
  });

  beforeEach(() => {
    cy.visit('http://localhost:3000/');
  });

  it('should edit a test', () => {
    cy.get(`[data-cy=test-actions-button-${testId}]`).click();
    cy.get('[data-cy=test-card-edit]').click();

    cy.get('[data-cy=edit-test-modal]').should('be.visible');
    cy.get('[data-cy=name]').clear().type('Edited Test');

    cy.get('[data-cy=edit-test-submit]').click();
    cy.get('[data-cy=test-details-name]').should('have.text', `Edited Test (v2)`);

    cy.location('pathname').should('match', /\/test\/.*/i, {timeout: 10000});
  });

  it('should edit a test from the test details', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);

    cy.get(`[data-cy=test-actions-button-${testId}]`).click();
    cy.get('[data-cy=test-card-edit]').click();

    cy.get('[data-cy=edit-test-modal]').should('be.visible');
    cy.get('[data-cy=name]').clear().type('Edited Test');

    cy.get('[data-cy=edit-test-submit]').click();
    cy.get('[data-cy=test-details-name]').should('have.text', `Edited Test (v2)`);

    cy.location('pathname').should('match', /\/run\/.*/i, {timeout: 10000});
  });
});
