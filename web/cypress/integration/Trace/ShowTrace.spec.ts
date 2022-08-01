describe('Show Trace', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should show the trace components', () => {
    cy.location('pathname').then(pathname => {
      cy.goToTestDetailPageAndRunTest(pathname);

      cy.get('[data-cy^=trace-node-]', {timeout: 30000}).should('be.visible');
      cy.get('[data-cy=span-details-attributes]').should('be.visible');
      cy.get('[data-cy=empty-assertion-card-list]').should('exist');

      cy.get('[data-cy=assertion-card-list]').should('exist');
    });
  });
});
