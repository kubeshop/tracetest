describe('Run Test', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should show and click the Run Test button when the test has finished', () => {
    cy.location('pathname').then(pathname => {
      cy.goToTestDetailPageAndRunTest(pathname);
      cy.get('[data-cy=run-test-button]', {timeout: 20000}).should('be.visible');
      cy.get(`[data-cy^=run-test-button]`).first().click();
    });
  });
});
