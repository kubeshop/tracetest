import TransactionUtils from '../utils/Transactions';

const transactionUtils = TransactionUtils();

describe('Transactions', () => {
  beforeEach(() => {
    cy.visit('/');
    cy.wrap(transactionUtils.createTests());
  });

  afterEach(() => {
    cy.wrap(transactionUtils.deleteTests()).then(() => {
      cy.deleteTransaction();
    });
  });

  it('should create and run a transaction with multiple tests', () => {
    const name = `Transaction - #${String(Date.now()).slice(-4)}`;
    cy.openTransactionCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, name, 'CreateTransactionFactory');

    transactionUtils.testList.forEach(test => {
      cy.get('[data-cy=transaction-test-selection]').click();
      cy.get(`[data-cy="${test.name}"]`).first().click();
    });

    cy.submitCreateForm('CreateTransactionFactory');
    cy.get('[data-cy=transaction-details-name').should('have.text', `${name} (v1)`);
    cy.get('[data-cy=transaction-run-button]', {timeout: 30000}).should('be.visible');
  });

  it('should rerun a transaction after creation', () => {
    const name = `Transaction - #${String(Date.now()).slice(-4)}`;
    cy.openTransactionCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, name, 'CreateTransactionFactory');

    transactionUtils.testList.forEach(test => {
      cy.get('[data-cy=transaction-test-selection]').click();
      cy.get(`[data-cy="${test.name}"]`).first().click();
    });

    cy.submitCreateForm('CreateTransactionFactory');
    cy.get('[data-cy=transaction-details-name').should('have.text', `${name} (v1)`);
    cy.get('[data-cy=transaction-run-button]', {timeout: 30000}).should('be.visible').click();

    cy.get('[data-cy=transaction-details-name').should('have.text', `${name} (v1)`);
    cy.get('[data-cy=transaction-run-button]', {timeout: 30000}).should('be.visible').click();
  });

  it('should update a transaction and rerun', () => {
    const name = `Transaction - #${String(Date.now()).slice(-4)}`;
    cy.openTransactionCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, name, 'CreateTransactionFactory');

    transactionUtils.testList.forEach(test => {
      cy.get('[data-cy=transaction-test-selection]').click();
      cy.get(`[data-cy="${test.name}"]`).first().click();
    });

    cy.submitCreateForm('CreateTransactionFactory');
    cy.get('[data-cy=transaction-details-name').should('have.text', `${name} (v1)`);
    cy.get('[data-cy=transaction-run-button]', {timeout: 50000}).should('be.visible');
    cy.get('[data-cy^=transaction-execution-step-]').should('have.length', 2);

    const updateName = `${name} - updated`;

    cy.get('[data-cy=create-test-name-input]').type(' - updated');
    transactionUtils.testList.forEach(test => {
      cy.get('[data-cy=transaction-test-selection]').click();
      cy.get(`[data-cy="${test.name}"]`).first().click();
    });

    cy.get('[data-cy=edit-transaction-submit]').click();
    cy.get('[data-cy=transaction-details-name').should('have.text', `${updateName} (v2)`);
    cy.get('[data-cy=transaction-run-button]', {timeout: 50000}).should('be.visible');
    cy.get('[data-cy^=transaction-execution-step-]').should('have.length', 4);
  });
});
