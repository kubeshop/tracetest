describe('Cancel Create test', () => {
  beforeEach(() => cy.visit('/'));

  it('should cancel a create test flow', () => {
    cy.navigateToTestCreationPage();
    cy.get('[data-cy=create-test-cancel]').click();
  });
});
