describe('Create test from Postman Collection', () => {
  beforeEach(() => {
    cy.enableDemo();
    cy.visit('/');
  });

  it('should create a basic GET test', () => {
    cy.interceptHomeApiCall();
    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    cy.openTestCreationModal();
    cy.get('[data-cy=postman-plugin]').click();
    cy.fillCreateFormBasicStep(name, 'Create from Postman Collection');
    cy.get('[data-cy="collectionFile"]').attachFile('collection.json');
    cy.get('[data-cy=collectionTest-select]').click();
    cy.get('[data-cy=collectionTest-1]').click({force: true});
    cy.submitCreateForm();
    cy.matchTestRunPageUrl();
    cy.cancelOnBoarding();
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
    cy.deleteTest(true);
  });
});
