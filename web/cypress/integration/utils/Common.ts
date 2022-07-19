import {camelCase} from 'lodash';
import {DemoTestExampleList} from '../../../src/constants/Test.constants';

Cypress.on('uncaught:exception', err => !err.message.includes('ResizeObserver loop limit exceeded'));

export const [{name, description}] = DemoTestExampleList;

// eslint-disable-next-line import/no-mutable-exports
export let testId = '';

export const createTest = () => {
  cy.visit('http://localhost:3000/');
  const $form = navigateToTestCreationPage();

  $form.get('[data-cy=create-test-next-button]').click();

  $form.get('[data-cy=example-button]').click();
  $form.get(`[data-cy=demo-example-${camelCase(name)}]`).click();

  $form.get('[data-cy=create-test-next-button]').last().click();
  $form.get('[data-cy=create-test-create-button]').last().click();

  cy.location('pathname').should('match', /\/test\/.*/i, {timeout: 20000});
  cy.location().then(({pathname}) => {
    const id = getTestId(pathname);

    testId = id;
  });
  cy.visit('http://localhost:3000/');
};

export const navigateToTestCreationPage = () => {
  cy.get('[data-cy=create-test-button]').click();

  const $form = cy.get('[data-cy=create-test-header]');
  $form.should('be.visible');

  return $form;
};

export const deleteTest = () => {
  cy.location().then(({pathname}) => {
    const localTestId = getTestId(pathname);
    cy.visit('http://localhost:3000/');

    cy.get(`[data-cy=test-actions-button-${localTestId}]`).click();
    cy.get('[data-cy=test-card-delete]').click();

    cy.get(`[data-cy=test-actions-button-${localTestId}]`).should('not.exist');
  });
};

export const getTestId = (pathname: string) => {
  const [, , localTestId] = pathname.split('/').reverse();

  return localTestId;
};

export const getResultId = (pathname: string) => {
  const [resultId, ,] = pathname.split('/').reverse();

  return resultId;
};

export const createMultipleTestRuns = (id: string, count: number) => {
  cy.visit('http://localhost:3000/');

  for (let i = 0; i < count; i += 1) {
    cy.get(`[data-cy=test-run-button-${id}]`).click();
    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.wait(500);

    cy.visit('http://localhost:3000/');
  }
};
