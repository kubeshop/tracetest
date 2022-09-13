import 'cypress-file-upload';
import {camelCase} from 'lodash';
import {Plugins} from '../../src/constants/Plugins.constants';
import {getTestId} from '../e2e/utils/Common';

export const testRunPageRegex = /\/test\/(.*)\/run\/(.*)/;
export const getAttributeListId = (number: number) => `#assertion-form_assertions_${number}_attribute_list`;
export const getComparatorListId = (number: number) => `#assertion-form_assertions_${number}_comparator_list`;

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
    cy.visit(`/`);
    cy.wait('@testList');
    cy.get('[data-cy=test-list]').should('exist', {timeout: 10000});
    cy.get(`[data-cy=test-actions-button-${localTestId}]`, {timeout: 10000}).should('be.visible');
    cy.get(`[data-cy=test-actions-button-${localTestId}]`).click({force: true});
    cy.get('[data-cy=test-card-delete]').click();
    cy.get('[data-cy=delete-confirmation-modal] .ant-btn-primary').click();
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

Cypress.Commands.add('interceptEditTestCall', () => {
  cy.intercept({method: 'PUT', url: '/api/tests/*'}).as('testEdit');
});

Cypress.Commands.add('inteceptHomeApiCall', () => {
  cy.intercept({method: 'GET', url: '/api/tests?take=20&skip=0&query='}).as('testList');
  cy.intercept({method: 'DELETE', url: '/api/tests/**'}).as('testDelete');
  cy.intercept({method: 'POST', url: '/api/tests'}).as('testCreation');
});

Cypress.Commands.add('waitForTracePageApiCalls', () => {
  cy.wait('@testRun');
  cy.wait('@testObject');
  // traces take some time to return
  cy.wait('@testRuns', {timeout: 30000});
});

Cypress.Commands.add('createTestWithAuth', (authMethod: string, keys: string[]): any => {
  cy.get('[data-cy=create-test-next-button]').last().click();
  cy.selectTestFromDemoList();
  cy.get('[data-cy=auth-type-select]').click();
  cy.get(`[data-cy=auth-type-select-option-${authMethod}]`).click();
  keys.forEach(key => cy.get(`[data-cy=${authMethod}-${key}]`).type(key));
  return cy.wrap(Plugins.REST.demoList[0].name);
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
  cy.visit(`/test/${testId}`);
  cy.get('[data-cy^=result-card]', {timeout: 10000}).first().click();
  cy.makeSureUserIsOnTestDetailPage();
  cy.get(`[data-cy^=test-run-result-]`).first().click();
  cy.makeSureUserIsOnTracePage(false);
});

Cypress.Commands.add('makeSureUserIsOnTestDetailPage', () => {
  cy.location('href').should('match', /\/test\/.*/i);
  cy.wait('@testObject');
});

Cypress.Commands.add('makeSureUserIsOnTracePage', (shouldCancelOnboarding = true) => {
  cy.matchTestRunPageUrl();
  if (shouldCancelOnboarding) {
    cy.cancelOnBoarding();
  }
});

Cypress.Commands.add('cancelOnBoarding', () => {
  const value = localStorage.getItem('guided_tour');
  const parsedValue = value ? JSON.parse(value) : undefined;

  if (!parsedValue || parsedValue.trace === false) {
    cy.get('[data-cy=no-thanks]').click();
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

Cypress.Commands.add('editTestByTestId', () => {
  cy.interceptEditTestCall();
  cy.get('[data-cy=edit-test-form]').should('be.visible');
  cy.get('[data-cy=create-test-name-input] input').clear().type('Edited Test');
  cy.get('[data-cy=edit-test-submit]').click();
  cy.wait('@testEdit');
  cy.wait('@testObject');
  cy.wait('@testRun');
  cy.get('[data-cy=test-details-name]').should('have.text', `Edited Test (v2)`);
  cy.matchTestRunPageUrl();
});

Cypress.Commands.add('selectOperator', (index: number, text?: string) => {
  cy.get('[data-cy=assertion-check-operator]').last().click();
  cy.get(`${getComparatorListId(index)} + div .ant-select-item`)
    .last()
    .click();
  if (text) {
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').last().should('have.text', text);
  }
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
  cy.cancelOnBoarding();
});

Cypress.Commands.add('createAssertion', (index = 0) => {
  cy.selectRunDetailMode(3);

  cy.get(`[data-cy=trace-node-database]`, {timeout: 25000}).first().click({force: true});
  cy.get('[data-cy=add-test-spec-button]').click({force: true});
  cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
  cy.get('[data-cy=assertion-check-attribute]').type('db');
  const attributeListId = getAttributeListId(index);
  cy.get(`${attributeListId} + div .ant-select-item:nth-child(2)`).first().click({force: true});
  cy.get('[data-cy=assertion-check-operator]').click({force: true});
  cy.get('[data-cy=assertion-form-submit-button]').click();
  cy.get('[data-cy=test-specs-container]').should('be.visible');
  cy.get('[data-cy=test-spec-container]').should('have.lengthOf', 1);
});

/**
 * Click the test run detail mode tabs
 * index: 1 = trigger, 2 = trace, 3 = test
 */
Cypress.Commands.add('selectRunDetailMode', (index: number) => {
  cy.get(`[data-cy=run-detail-header] .ant-tabs-nav-list div:nth-child(${index})`).click();
});
