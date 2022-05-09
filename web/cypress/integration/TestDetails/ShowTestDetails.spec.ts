import {camelCase} from 'lodash';
import {openCreateTestModal, deleteTest} from '../Home/CreateTest.spec';
import {DemoTestExampleList} from '../../../src/constants/Test.constants';

const [, {name, description}] = DemoTestExampleList;
let testId = '';

export const createTest = () => {
  cy.visit('http://localhost:3000/');
  const $form = openCreateTestModal();

  $form.get('[data-cy=example-button]').click();
  $form.get(`[data-cy=demo-example-${camelCase(name)}]`).click();
  $form.get('[data-cy=create-test-submit]').click();

  cy.location('pathname').should('match', /\/test\/.*/i);
  cy.location().then(({pathname}) => {
    testId = pathname.split('/').pop();
  });
  cy.visit('http://localhost:3000/');
};

describe('Show test details', () => {
  before(() => {
    createTest();
  });

  after(() => {
    deleteTest();
  });

  it('should show the test details for any trace', () => {
    cy.get(`[data-cy=test-url-${testId}]`).click();

    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} - ${description}`);
    cy.get('[data-cy=test-result-table]').should('be.visible');
    cy.get('[data-cy=test-result-table] .ant-table-row').should('have.length.above', 0);
  });

  it('should run a new test', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);
    let testRunResultId = '';
    cy.get(`[data-cy=test-details-run-test-button]`).click();
    cy.location('href').should('match', /resultId=.*/i);

    cy.location().then(({href}) => {
      testRunResultId = href.split('?').pop().split('=').pop();

      cy.wait(2000);
      cy.get('#rc-tabs-0-tab-1').click();
      cy.get(`[data-cy=test-run-result-${testRunResultId}]`).should('be.visible');
    });
  });

  it('should update the test run result status', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);
    let testRunResultId = '';
    cy.get(`[data-cy=test-details-run-test-button]`).click();
    cy.location('href').should('match', /resultId=.*/i);

    cy.location().then(({href}) => {
      testRunResultId = href.split('?').pop().split('=').pop();

      cy.wait(2000);
      cy.get('[data-cy=test-run-result-status]').should('have.text', 'Test status:Awaiting trace');
      cy.wait(2000);
      cy.get('[data-cy=test-run-result-status]').should('have.text', 'Test status:Finished');
      cy.get('#rc-tabs-0-tab-1').click();
      cy.get(`[data-cy=test-run-result-status-${testRunResultId}]`).should('have.text', 'Finished');
    });
  });
});
