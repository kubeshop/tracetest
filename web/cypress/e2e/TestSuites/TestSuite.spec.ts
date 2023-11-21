import TestSuite from '../utils/TestSuites';

const testSuiteUtils = TestSuite();

describe('TestSuites', () => {
  beforeEach(() => {
    cy.visit('/testsuites');
    cy.wrap(testSuiteUtils.createTests());
  });

  afterEach(() => {
    cy.wrap(testSuiteUtils.deleteTests()).then(() => {
      cy.deleteTestSuite();
    });
  });

  it('should create a TestSuite with multiple tests', () => {
    const name = `TestSuite - #${String(Date.now()).slice(-4)}`;
    cy.openTestSuiteCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, 'CreateTestSuiteFactory');

    testSuiteUtils.testList.forEach(test => {
      cy.get('[data-cy=testsuite-test-selection]').click();
      cy.get(`[data-cy="${test.spec.name}"]`).first().click();
    });

    cy.submitCreateForm('CreateTestSuiteFactory');
    cy.get('[data-cy=testsuite-details-name').should('have.text', `${name} (v1)`);
  });

  it('should create a TestSuite with no tests', () => {
    const name = `TestSuite - #${String(Date.now()).slice(-4)}`;
    cy.openTestSuiteCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, 'CreateTestSuiteFactory');

    cy.submitCreateForm('CreateTestSuiteFactory');
    cy.get('[data-cy=testsuite-details-name').should('have.text', `${name} (v1)`);
  });
});
