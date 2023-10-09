describe('Show Onboarding', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should init the onboarding flow', () => {
    cy.get('[data-cy=menu-link]').click();
    cy.get('[data-cy=menu-onboarding]').click();

    cy.get('[data-cy=onboarding-container]').within(() => {
      cy.get('[data-cy=onboarding-step]').should('contain.text', '1 of 5');
      cy.get('[data-cy=onboarding-next]').click();

      cy.get('[data-cy=onboarding-step]').should('contain.text', '2 of 5');
      cy.get('[data-cy=onboarding-next]').click();

      cy.get('[data-cy=onboarding-step]').should('contain.text', '3 of 5');
      cy.get('[data-cy=onboarding-next]').click();

      cy.get('[data-cy=onboarding-step]').should('contain.text', '4 of 5');
      cy.get('[data-cy=onboarding-next]').click();

      cy.get('[data-cy=onboarding-step]').should('contain.text', '5 of 5');
      cy.get('[data-cy=onboarding-next]').click();
    });
  });
});
