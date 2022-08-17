import 'cypress-file-upload';
import {camelCase} from 'lodash';
import {Plugins} from '../../src/constants/Plugins.constants';
import {getTestId, testRunPageRegex} from '../integration/utils/Common';

Cypress.Commands.add('createMultipleTestRuns', (id: string, count: number) => {
  cy.visit('/');

  for (let i = 0; i < count; i += 1) {
    cy.get(`[data-cy=test-run-button-${id}]`).click();
    cy.matchTestRunPageUrl();

    cy.visit('/');
  }
});

Cypress.Commands.add('setCreateFormUrl', (method: string, url: string) => {
  cy.get('[data-cy=method-select]').click();
  cy.get(`[data-cy=method-select-option-${method}]`).click();
  cy.get('[data-cy=url]').type(url);
});

Cypress.Commands.add('deleteTest', (shoudlIntercept = false) => {
  cy.location('pathname').then(pathname => {
    const localTestId = getTestId(pathname);
    // called when test not created with createTest method
    if (shoudlIntercept) {
      cy.inteceptHomeApiCall();
    }
    cy.visit(`http://localhost:3000`);
    cy.wait('@testList');
    cy.get(`[data-cy=test-actions-button-${localTestId}]`).should('be.visible');
    cy.get(`[data-cy=test-actions-button-${localTestId}]`).click({force: true});
    cy.get('[data-cy=test-card-delete]').click();
    cy.wait('@testDelete');
    cy.get(`[data-cy=test-actions-button-${localTestId}]`).should('not.exist');
    cy.wait('@testList');
  });
});

Cypress.Commands.add('navigateToTestCreationPage', () => {
  cy.get('[data-cy=create-test-button]').click();
  cy.get('[data-cy=create-test-header]').should('be.visible');
});

Cypress.Commands.add('interceptTracePageApiCalls', () => {
  cy.intercept({method: 'GET', url: '/api/tests/**/run/**'}).as('testRun');
  cy.intercept({method: 'GET', url: '/api/tests/**'}).as('testObject');
  cy.intercept({method: 'PUT', url: '/api/tests/**/run/**/dry-run'}).as('testRuns');
});

Cypress.Commands.add('inteceptHomeApiCall', () => {
  cy.intercept({method: 'GET', url: '/api/tests?take=20&skip=0&query='}).as('testList');
  cy.intercept({method: 'DELETE', url: '/api/tests/**'}).as('testDelete');
  cy.intercept({method: 'POST', url: '/api/tests'}).as('testCreation');
});

Cypress.Commands.add('waitForTracePageApiCalls', () => {
  cy.wait('@testRun');
  cy.wait('@testObject');
  cy.wait('@testRuns');
});

Cypress.Commands.add('createTestWithAuth', (authMethod: string, keys: string[]) => {
  const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
  cy.fillCreateFormBasicStep(name);
  cy.setCreateFormUrl('GET', 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon');
  cy.get('[data-cy=auth-type-select]').click();
  cy.get(`[data-cy=auth-type-select-option-${authMethod}]`).click();
  keys.forEach(key => cy.get(`[data-cy=${authMethod}-${key}]`).type(key));
  return cy.wrap(name);
});

Cypress.Commands.add('submitAndMakeSureTestIsCreated', (name: string) => {
  cy.submitCreateTestForm();
  cy.interceptTracePageApiCalls();
  cy.makeSureUserIsOnTracePage();
  cy.waitForTracePageApiCalls();
  cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
  cy.deleteTest(true);
});

Cypress.Commands.add('matchTestRunPageUrl', () => {
  cy.location('pathname').should('match', testRunPageRegex);
});

Cypress.Commands.add('goToTestDetailPageAndRunTest', (pathname: string) => {
  const testId = getTestId(pathname);
  cy.visit(`http://localhost:3000/test/${testId}`);
  cy.get('[data-cy^=result-card]', {timeout: 10000}).first().click();
  cy.makeSureUserIsOnTestDetailPage();
  cy.get(`[data-cy^=test-run-result-]`).first().click();
  cy.makeSureUserIsOnTracePage(false);
});

Cypress.Commands.add('makeSureUserIsOnTestDetailPage', () => {
  cy.location('href').should('match', /\/test\/.*/i);
});

Cypress.Commands.add('makeSureUserIsOnTracePage', (shouldCancelOnboarding = true) => {
  cy.matchTestRunPageUrl();
  if (shouldCancelOnboarding) {
    cy.cancelOnBoarding();
  }
});

Cypress.Commands.add('cancelOnBoarding', () => {
  if (cy.get('[data-cy=no-thanks]', {timeout: 10000})) {
    cy.get('[data-cy=no-thanks]', {timeout: 10000}).click();
  }
});

Cypress.Commands.add('submitCreateTestForm', () => {
  cy.get('[data-cy=create-test-create-button]').last().click();
  cy.wait('@testCreation');
});

Cypress.Commands.add('fillCreateFormBasicStep', (name: string, description?: string) => {
  cy.get('[data-cy=create-test-next-button]').click();
  cy.get('[data-cy=create-test-name-input').type(name);
  cy.get('[data-cy=create-test-description-input').type(description || name);
  cy.get('[data-cy=create-test-next-button]').last().click();
});

Cypress.Commands.add('createTestByName', (name: string) => {
  cy.navigateToTestCreationPage();
  cy.get('[data-cy=create-test-next-button]').click();
  cy.get('[data-cy=example-button]').click();
  cy.get(`[data-cy=demo-example-${camelCase(name)}]`).click();
  cy.get('[data-cy=create-test-next-button]').last().click();
  cy.submitCreateTestForm();
  cy.makeSureUserIsOnTracePage();
  cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
});

Cypress.Commands.add('editTestByTestId', (testId: string) => {
  cy.get(`[data-cy=test-actions-button-${testId}]`).click();
  cy.get('[data-cy=test-card-edit]').click();
  cy.get('[data-cy=edit-test-form]').should('be.visible');
  cy.get('[data-cy=create-test-name-input] input').clear().type('Edited Test');
  cy.get('[data-cy=edit-test-submit]').click();
  cy.get('[data-cy=test-details-name]').should('have.text', `Edited Test (v2)`);
  cy.matchTestRunPageUrl();
});

Cypress.Commands.add('selectTestFromDemoList', () => {
  cy.get('[data-cy=example-button]').click();
  cy.get(`[data-cy=demo-example-${camelCase(Plugins.REST.demoList[0].name)}]`).click();
  cy.get('[data-cy=create-test-next-button]').last().click();
});

Cypress.Commands.add('clickNextOnCreateTestWizard', () => {
  cy.get('[data-cy=create-test-next-button]').click();
});

Cypress.Commands.add('createTest', () => {
  cy.inteceptHomeApiCall();
  cy.visit('/');
  cy.navigateToTestCreationPage();
  cy.clickNextOnCreateTestWizard();
  cy.selectTestFromDemoList();
  cy.interceptTracePageApiCalls();
  cy.submitCreateTestForm();
  cy.makeSureUserIsOnTracePage();
  cy.waitForTracePageApiCalls();
});
