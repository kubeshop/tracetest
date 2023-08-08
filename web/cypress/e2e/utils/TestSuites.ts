import {TRawTestResource} from '../../../src/models/Test.model';
import {testSuiteTestList} from '../constants/TestSuite';

interface ITestSuiteUtils {
  testList: TRawTestResource[];
  createTest(test: TRawTestResource): Promise<TRawTestResource>;
  createTests(tests?: TRawTestResource[]): Promise<TRawTestResource[]>;
  deleteTest(id: string): Promise<void>;
  deleteTests(): Promise<void[]>;
  waitForTestSuiteRun(): void;
}

const TestSuitesUtils = (): ITestSuiteUtils => ({
  testList: [],
  createTest(test) {
    return new Promise(resolve => {
      cy.request('POST', '/api/tests', test).then((res: Cypress.Response<TRawTestResource>) => {
        resolve(res.body);
      });
    });
  },
  async createTests(tests = testSuiteTestList) {
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
  waitForTestSuiteRun() {
    cy.get('[data-cy=testsuite-run-result-status]').should('have.text', 'FINISHED', {timeout: 60000});
  }
});

export default TestSuitesUtils;
