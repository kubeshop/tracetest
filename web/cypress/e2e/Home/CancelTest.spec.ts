describe('Cancel Create test', () => {
  beforeEach(() => cy.visit('/'));

  it('should cancel a create test flow', () => {
    cy.openTestCreationModal();
    cy.get('.ant-modal-close-icon').click();
  });
});
