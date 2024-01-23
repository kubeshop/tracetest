describe('Home', () => {
  beforeEach(() => cy.visit('/tests'));

  it('should render the layout', () => {
    cy.get('[data-cy=menu-link]').should('be.visible');
  });

  it('should render the list of tests', () => {
    cy.get('[data-cy=create-button]').should('be.visible');
    cy.get('[data-cy=test-list]').should('exist');
  });

  it('should run a test from the home page', () => {
    cy.createTest();
    cy.visit('/tests');
    cy.get('[data-cy^=test-run-button]:not([data-cy*=button-00])', {timeout: 10000}).first().click();
    cy.location('href').should('match', /\/test\/.*/i);
    cy.deleteTest();
  });
});
