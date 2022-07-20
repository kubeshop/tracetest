import {camelCase} from 'lodash';
import Bluebird = require('cypress/types/bluebird');
import {Plugins} from '../../../src/constants/Plugins.constants';

Cypress.on('uncaught:exception', err => !err.message.includes('ResizeObserver loop limit exceeded'));

export const [{name, description}] = Plugins.REST.demoList;

export const testRunPageRegex = /\/test\/(.*)\/run\/(.*)/;

function promisify(chain): Bluebird<unknown> {
  return new Cypress.Promise((resolve, reject) => {
    // We must subscribe to failures and bail. Without this, the Cypress runner would never stop
    Cypress.on('fail', rejectPromise);
    // // unsubscribe from test failure on both success and failure. This cleanup is essential
    function resolvePromise(value) {
      resolve(value);
      Cypress.off('fail', rejectPromise);
    }

    function rejectPromise(error) {
      reject(error);
      Cypress.off('fail', rejectPromise);
    }

    chain.then(resolvePromise);
  });
}

export function interceptTracePageApiCalls() {
  cy.intercept({method: 'GET', url: '/api/tests/**/run/**'}).as('testRun');
  cy.intercept({method: 'GET', url: '/api/tests/**'}).as('testObject');
  cy.intercept({method: 'PUT', url: '/api/tests/**/run/**/dry-run'}).as('testRuns');
}

export function inteceptHomeApiCall() {
  cy.intercept({method: 'GET', url: '/api/tests'}).as('testList');
}

export function waitForTracePageApiCalls() {
  cy.wait('@testRun');
  cy.wait('@testObject');
  cy.wait('@testRuns', {timeout: 20000});
}

export function makeSureUserisOnTracePage() {
  cy.location('pathname').should('match', testRunPageRegex, {timeout: 20000}).wait(2000);
}

export async function extractTestIdFromTracePage(): Promise<string> {
  const pathname = (await promisify(cy.location('pathname'))) as string;
  return getTestId(pathname);
}

export const createTest = async (): Promise<string> => {
  cy.server();
  inteceptHomeApiCall();
  cy.visit('http://localhost:3000/');
  const $form = navigateToTestCreationPage();
  $form.get('[data-cy=create-test-next-button]').click();
  $form.get('[data-cy=example-button]').click();
  $form.get(`[data-cy=demo-example-${camelCase(name)}]`).click();
  $form.get('[data-cy=create-test-next-button]').last().click();
  interceptTracePageApiCalls();
  $form.get('[data-cy=create-test-create-button]').last().click();
  makeSureUserisOnTracePage();
  cy.get('[data-cy=no-thanks').click();
  waitForTracePageApiCalls();
  const testId = await extractTestIdFromTracePage();
  cy.log(testId);
  return testId;
};

export const navigateToTestCreationPage = () => {
  cy.get('[data-cy=create-test-button]').click();

  const $form = cy.get('[data-cy=create-test-header]');
  $form.should('be.visible');

  return $form;
};

function deleteTestByID(localTestId: string) {
  cy.get(`[data-cy=test-actions-button-${localTestId}]`).should('be.visible');
  cy.get(`[data-cy=test-actions-button-${localTestId}]`).click({force: true});
  cy.get('[data-cy=test-card-delete]').click();
  cy.get(`[data-cy=test-actions-button-${localTestId}]`).should('not.exist');
}

export const deleteTest = (testId?: string) => {
  cy.visit(`http://localhost:3000`);
  cy.wait('@testList').wait(2000);
  deleteTestByID(testId);
};

export const getTestId = (pathname: string) => {
  const [, , localTestId] = pathname.split('/').reverse();

  return localTestId;
};

export const getResultId = (pathname: string) => {
  const [resultId] = pathname.split('/').reverse();

  return resultId;
};

export const createMultipleTestRuns = (id: string, count: number) => {
  cy.visit('http://localhost:3000/');

  for (let i = 0; i < count; i += 1) {
    cy.get(`[data-cy=test-run-button-${id}]`).click();
    cy.location('pathname').should('match', testRunPageRegex);
    cy.wait(500);

    cy.visit('http://localhost:3000/');
  }
};
