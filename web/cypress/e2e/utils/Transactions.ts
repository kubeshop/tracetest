import {TRawTestResource} from '../../../src/models/Test.model';
import {transactionTestList} from '../constants/Transactions';

interface ITransactionUtils {
  testList: TRawTestResource[];
  createTest(test: TRawTestResource): Promise<TRawTestResource>;
  createTests(tests?: TRawTestResource[]): Promise<TRawTestResource[]>;
  deleteTest(id: string): Promise<void>;
  deleteTests(): Promise<void[]>;
  waitForTransactionRun(): void;
}

const TransactionUtils = (): ITransactionUtils => ({
  testList: [],
  createTest(test) {
    return new Promise(resolve => {
      cy.request('POST', '/api/tests', test).then((res: Cypress.Response<TRawTestResource>) => {
        resolve(res.body);
      });
    });
  },
  async createTests(tests = transactionTestList) {
    this.testList = await Promise.all(tests.map(test => this.createTest(test)));

    return this.testList;
  },
  deleteTest(id) {
    return new Promise(resolve => {
      cy.request('DELETE', `/api/tests/${id}`).then(() => {
        resolve();
      });
    });
  },
  deleteTests() {
    return Promise.all(this.testList.map(test => this.deleteTest(test.spec.id)));
  },
  waitForTransactionRun() {
    cy.get('[data-cy=transaction-run-result-status]').should('have.text', 'FINISHED', {timeout: 60000});
  }
});

export default TransactionUtils;
