import TestSuites from '../utils/TestSuites';

const testSuitesUtils = TestSuites();

describe('TestSuiteRuns', () => {
  beforeEach(() => {
    cy.visit('/testsuites');
    cy.wrap(testSuitesUtils.createTests());
  });

  afterEach(() => {
    cy.wrap(testSuitesUtils.deleteTests()).then(() => {
      cy.deleteTestSuite();
    });
  });

  it('should create and run a test suite with multiple tests', () => {
    const name = `TestSuite - #${String(Date.now()).slice(-4)}`;
    cy.openTestSuiteCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, name, 'CreateTestSuiteFactory');

    testSuitesUtils.testList.forEach(test => {
      cy.get('[data-cy=testsuite-test-selection]').click();
      cy.get(`[data-cy="${test.spec.name}"]`).first().click();
    });

    cy.submitCreateForm('CreateTestSuiteFactory');
    cy.get('[data-cy=testsuite-details-name]').should('have.text', `${name} (v1)`);
    cy.reload().get('[data-cy=testsuite-run-button]', {timeout: 50000}).should('be.visible');
  });

  it('should rerun a test suite after creation', () => {
    const name = `TestSuite - #${String(Date.now()).slice(-4)}`;
    cy.openTestSuiteCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, name, 'CreateTestSuiteFactory');

    testSuitesUtils.testList.forEach(test => {
      cy.get('[data-cy=testsuite-test-selection]').click();
      cy.get(`[data-cy="${test.spec.name}"]`).first().click();
    });

    cy.submitCreateForm('CreateTestSuiteFactory');
    cy.get('[data-cy=testsuite-details-name').should('have.text', `${name} (v1)`);
    cy.reload().get('[data-cy=testsuite-run-button]', {timeout: 50000}).should('be.visible');

    cy.get('[data-cy=testsuite-details-name').should('have.text', `${name} (v1)`);
    cy.reload().get('[data-cy=testsuite-run-button]', {timeout: 50000}).should('be.visible');
  });

  it('should update a test suite and rerun', () => {
    const name = `TestSuite - #${String(Date.now()).slice(-4)}`;
    cy.openTestSuiteCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, name, 'CreateTestSuiteFactory');

    testSuitesUtils.testList.forEach(test => {
      cy.get('[data-cy=testsuite-test-selection]').click();
      cy.get(`[data-cy="${test.spec.name}"]`).first().click();
    });

    cy.submitCreateForm('CreateTestSuiteFactory');
    cy.get('[data-cy=testsuite-details-name').should('have.text', `${name} (v1)`);
    cy.reload().get('[data-cy=testsuite-run-button]', {timeout: 50000}).should('be.visible');
    cy.get('[data-cy^=testsuite-execution-step-]').should('have.length', 2);

    const updateName = `${name} - updated`;

    cy.get('[data-cy=create-test-name-input]').type(' - updated');
    testSuitesUtils.testList.forEach(test => {
      cy.get('[data-cy=testsuite-test-selection]').click();
      cy.get(`[data-cy="${test.spec.name}"]`).first().click();
    });

    cy.get('[data-cy=edit-testsuite-submit]').click();
    cy.get('[data-cy=testsuite-details-name').should('have.text', `${updateName} (v2)`);
    cy.reload().get('[data-cy=testsuite-run-button]', {timeout: 50000}).should('be.visible');
    cy.get('[data-cy^=testsuite-execution-step-]').should('have.length', 4);
  });
});
