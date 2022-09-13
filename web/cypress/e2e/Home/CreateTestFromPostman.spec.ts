describe('Create test from Postman Collection', () => {
  beforeEach(() => cy.visit('/'));

  it('should create a basic GET test', () => {
    cy.inteceptHomeApiCall();
    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    cy.navigateToTestCreationPage();
    cy.get('[data-cy=postman-plugin]').click();
    cy.fillCreateFormBasicStep(name, 'Create from Postman Collection');
    cy.get('[data-cy="collectionFile"]').attachFile('collection.json');
    cy.get('[data-cy=collectionTest-select]').click();
    cy.get('[data-cy=collectionTest-1]').click({force: true});
    cy.submitCreateTestForm();
    cy.matchTestRunPageUrl();
    cy.cancelOnBoarding();
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
    cy.deleteTest(true);
  });
});
