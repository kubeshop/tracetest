import 'cypress-file-upload';
import {camelCase} from 'lodash';
import {POKEMON_HTTP_ENDPOINT, PokeshopDemo} from '../e2e/constants/Test';
import {getTestId, getTestSuiteId} from '../e2e/utils/Common';
import { SupportedPlugins } from '../../src/constants/Common.constants';

export const testRunPageRegex = /\/test\/(.*)\/run\/(.*)/;
export const getAttributeListId = (number: number) => `.cm-tooltip-autocomplete [id$=-${number}]`;
export const getComparatorListId = (number: number) => `#assertion-form_assertions_${number}_comparator_list`;
export const getValueFromList = (number: number) => `.cm-tooltip-autocomplete li:nth-child(${number})`;

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

Cypress.Commands.add('deleteTest', (shouldIntercept = false) => {
  cy.location('pathname').then(pathname => {
    const localTestId = getTestId(pathname);
    // called when test not created with createTest method
    if (shouldIntercept) {
      cy.interceptHomeApiCall();
    }
    cy.visit(`/`);
    cy.wait('@testList');
    cy.get('[data-cy=test-list]').should('exist', {timeout: 10000});
    cy.get(`[data-cy=test-actions-button-${localTestId}]`, {timeout: 10000}).should('be.visible');
    cy.get(`[data-cy=test-actions-button-${localTestId}]`).click({force: true});
    cy.get('[data-cy=test-card-delete]').click();
    cy.get('[data-cy=confirmation-modal] .ant-btn-primary').click();
    cy.wait('@testDelete');
    cy.get(`[data-cy=test-actions-button-${localTestId}]`).should('not.exist');
    cy.wait('@testList');
    cy.clearLocalStorage();
  });
});

Cypress.Commands.add('openTestCreationModal', () => {
  cy.get('[data-cy=create-button]').click();
});

Cypress.Commands.add('interceptTracePageApiCalls', () => {
  cy.intercept({method: 'GET', url: /\/api\/tests\/([\w-]+)\/run\/(\w+)$/}).as('testRun');
  cy.intercept({method: 'GET', url: /\/api\/tests\/([\w-]+)$/}).as('testObject');
  cy.intercept({method: 'PUT', url: '/api/tests/**/run/**/dry-run'}).as('testRuns');
});

Cypress.Commands.add('interceptEditTestCall', () => {
  cy.intercept({method: 'PUT', url: '/api/tests/*'}).as('testEdit');
});

Cypress.Commands.add('interceptHomeApiCall', () => {
  cy.intercept({method: 'GET', url: '/api/tests?take=20&skip=0*'}).as('testList');
  cy.intercept({method: 'GET', url: '/api/testsuites?take=20&skip=0*'}).as('testSuitesList');
  cy.intercept({method: 'DELETE', url: '/api/tests/**'}).as('testDelete');
  cy.intercept({method: 'POST', url: '/api/tests'}).as('testCreation');
  cy.intercept({method: 'DELETE', url: '/api/testsuites/**'}).as('testSuiteDelete');
  cy.intercept({method: 'POST', url: '/api/testsuites'}).as('testSuiteCreation');
});

Cypress.Commands.add('waitForTracePageApiCalls', () => {
  cy.wait('@testRun');
  cy.wait('@testObject');
  // traces take some time to return
  cy.wait('@testRuns', {timeout: 30000});
});

Cypress.Commands.add('createTestWithAuth', (authMethod: string, keys: string[]): any => {
  cy.get('[data-cy=CreateTestFactory-create-next-button]').last().click();
  cy.selectTestFromDemoList();
  cy.get('[data-cy=auth-type-select]').click();
  cy.get(`[data-cy=auth-type-select-option-${authMethod}]`).click();
  keys.forEach(key => cy.get(`[data-cy=${authMethod}-${key}] [contenteditable]`).first().type('key'));
  return cy.wrap(PokeshopDemo[0].name);
});

Cypress.Commands.add('submitAndMakeSureTestIsCreated', (name: string) => {
  cy.submitCreateForm();
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
  cy.get('[data-cy^=run-card]', {timeout: 10000}).first().click();
  cy.makeSureUserIsOnTestDetailPage();
  cy.makeSureUserIsOnTracePage();
});

Cypress.Commands.add('makeSureUserIsOnTestDetailPage', () => {
  cy.location('href').should('match', /\/test\/.*/i);
  cy.wait('@testObject');
});

Cypress.Commands.add('makeSureUserIsOnTracePage', () => {
  cy.matchTestRunPageUrl();
  cy.cancelOnBoarding();
});

Cypress.Commands.add('cancelOnBoarding', () => {
  const value = localStorage.getItem('user_preferences');
  const parsedValue = value ? JSON.parse(value) : undefined;

  if (!parsedValue || parsedValue.showGuidedTourNotification === true) {
    cy.get('body').then($body => {
      if ($body.find('[data-cy=guided-tour-cancel-notification]').length > 0)
        cy.get('[data-cy=guided-tour-cancel-notification]').click();
    });
  }
});

Cypress.Commands.add('submitCreateForm', () => {
  cy.get(`[data-cy="run-test-submit"]`).click();
  cy.wait('@testCreation');
});

Cypress.Commands.add('fillCreateFormBasicStep', (name: string) => {
  cy.get('[data-cy=overlay-input-overlay]').should('be.visible').click();
  cy.get('[data-cy="overlay-input"]').clear().type(name);
});

Cypress.Commands.add('createTestByName', (name: string) => {
  cy.openTestCreationModal();
  cy.get(`[data-cy=${SupportedPlugins.REST.toLowerCase()}-plugin]`).click();
  cy.get('[data-cy=example-button]').click();
  cy.get(`[data-cy=demo-example-${camelCase(name)}]`).click();
  cy.submitCreateForm();
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
  cy.get(`[data-cy=demo-example-${camelCase(PokeshopDemo[0].name)}]`).click();
  cy.get('[data-cy=CreateTestFactory-create-next-button]').last().click();
});

Cypress.Commands.add('clickNextOnCreateTestWizard', () => {
  cy.get('[data-cy=CreateTestFactory-create-next-button]').click();
});

Cypress.Commands.add('createTest', () => {
  cy.enableDemo();
  cy.interceptHomeApiCall();
  cy.clearLocalStorage();
  cy.visit('/');
  cy.clearLocalStorage();
  cy.openTestCreationModal();
  cy.clickNextOnCreateTestWizard();
  cy.selectTestFromDemoList();
  cy.interceptTracePageApiCalls();
  cy.submitCreateForm();
  cy.makeSureUserIsOnTracePage();
  cy.waitForTracePageApiCalls();
});

Cypress.Commands.add('createAssertion', () => {
  cy.selectRunDetailMode(3);

  cy.get(`[data-cy=trace-node-database]`, {timeout: 25000}).first().click({force: true});
  cy.get('[data-cy=add-test-spec-button]').click({force: true});
  cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
  cy.get('[data-cy=editor-fallback]').should('not.exist');

  cy.get('[data-cy=expression-editor] [contenteditable]').first().type('db.name', {delay: 100});

  const attributeListId = getAttributeListId(0);
  cy.get(attributeListId, {timeout: 10000}).first().click();
  cy.get('[data-cy=assertion-check-value] .cm-content').first().click();
  cy.get(getValueFromList(1)).first().click();

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

Cypress.Commands.add('openTestSuiteCreationModal', () => {
  cy.get('[data-cy=create-button]').click();
  cy.get('[data-cy=create-test-steps-CreateTestSuiteFactory]').should('be.visible');
});

Cypress.Commands.add('deleteTestSuite', () => {
  cy.location('pathname').then(pathname => {
    const localTestId = getTestSuiteId(pathname);

    cy.visit(`/testsuites`);
    cy.wait('@testSuitesList');
    cy.get('[data-cy=test-list]').should('exist', {timeout: 10000});
    cy.get(`[data-cy=test-actions-button-${localTestId}]`, {timeout: 10000}).should('be.visible');
    cy.get(`[data-cy=test-actions-button-${localTestId}]`).click({force: true});
    cy.get('[data-cy=test-card-delete]').click();
    cy.get('[data-cy=confirmation-modal] .ant-btn-primary').click();
    cy.wait('@testSuiteDelete');
    cy.get(`[data-cy=test-actions-button-${localTestId}]`).should('not.exist');
    cy.wait('@testSuitesList');
    cy.clearLocalStorage();
  });
});

Cypress.Commands.add('enableDemo', () => {
  cy.visit(`/settings`);

  cy.get('[id*=tab-demo]').click();
  cy.get('#demo_pokeshop_enabled').then(element => {
    const isChecked = element.attr('aria-checked');

    if (isChecked === 'false') {
      cy.get('#demo_pokeshop_enabled').click();
      cy.get('#demo_pokeshop_pokeshop_httpEndpoint').type(POKEMON_HTTP_ENDPOINT);
      cy.get('[data-cy=demo-form-save-button]').click();
    }

    cy.visit(`/`);
  });
});
