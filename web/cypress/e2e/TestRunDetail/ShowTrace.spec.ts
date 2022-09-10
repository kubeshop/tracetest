describe('Show Trace', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should show the trace components', () => {
    cy.location('pathname').then(pathname => {
      cy.selectRunDetailMode(3);
      cy.get('[data-cy^=trace-node-]', {timeout: 30000}).should('be.visible');
      cy.get('[data-cy=empty-test-specs]').should('exist');
      cy.log('here');
      cy.goToTestDetailPageAndRunTest(pathname);
    });
  });
});
