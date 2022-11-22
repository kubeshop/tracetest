import {TRawTest} from '../../../src/types/Test.types';
import {transactionTestList} from '../constants/Transactions';

interface ITransactionUtils {
  testList: TRawTest[];
  createTest(test: TRawTest): Promise<TRawTest>;
  createTests(tests?: TRawTest[]): Promise<TRawTest[]>;
  deleteTest(id: string): Promise<void>;
  deleteTests(): Promise<void[]>;
  waitForTransactionRun(): void;
}

const TransactionUtils = (): ITransactionUtils => ({
  testList: [],
  createTest(test) {
    return new Promise(resolve => {
      cy.request('POST', '/api/tests', test).then((res: Cypress.Response<TRawTest>) => {
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
    return Promise.all(this.testList.map(test => this.deleteTest(test.id)));
  },
  waitForTransactionRun() {
    cy.get('[data-cy=transaction-run-result-status]').should('have.text', 'FINISHED', {timeout: 60000});
  }
});

export default TransactionUtils;
