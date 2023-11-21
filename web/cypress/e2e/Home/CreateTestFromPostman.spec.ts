describe('Create test from Postman Collection', () => {
  beforeEach(() => {
    cy.enableDemo();
    cy.visit('/');
  });

  it('should create a basic GET test', () => {
    cy.interceptHomeApiCall();
    cy.get('[data-cy=import-button]').click();
    cy.get('[data-cy=postman-plugin]').click();
    cy.get('[data-cy="collectionFile"]').attachFile('collection.json');
    cy.get('[data-cy=collectionTest-select]').click();
    cy.get('[data-cy=collectionTest-1]').click({force: true});
    cy.get(`[data-cy="import-test-submit"]`).click();
    cy.submitCreateForm();
    cy.matchTestRunPageUrl();
    cy.cancelOnBoarding();
    cy.get('[data-cy=overlay-input-overlay]').should('contain.text', 'create Test');
    cy.deleteTest(true);
  });
});
