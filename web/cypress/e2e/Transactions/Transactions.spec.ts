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

  it('should create a transaction with multiple tests', () => {
    const name = `Transaction - #${String(Date.now()).slice(-4)}`;
    cy.openTransactionCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, name, 'CreateTransactionFactory');

    transactionUtils.testList.forEach(test => {
      cy.get('[data-cy=transaction-test-selection]').click();
      cy.get(`[data-cy="${test.spec.name}"]`).first().click();
    });

    cy.submitCreateForm('CreateTransactionFactory');
    cy.get('[data-cy=transaction-details-name').should('have.text', `${name} (v1)`);
  });

  it('should create a transaction with no tests', () => {
    const name = `Transaction - #${String(Date.now()).slice(-4)}`;
    cy.openTransactionCreationModal();
    cy.interceptHomeApiCall();
    cy.fillCreateFormBasicStep(name, name, 'CreateTransactionFactory');

    cy.submitCreateForm('CreateTransactionFactory');
    cy.get('[data-cy=transaction-details-name').should('have.text', `${name} (v1)`);
  });
});
