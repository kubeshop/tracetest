import {camelCase} from 'lodash';
import {DemoTestExampleList} from '../../../src/constants/Test.constants';

Cypress.on('uncaught:exception', err => !err.message.includes('ResizeObserver loop limit exceeded'));

export const [, {name, description}] = DemoTestExampleList;

// eslint-disable-next-line import/no-mutable-exports
export let testId = '';

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

export const openCreateTestModal = () => {
  cy.get('[data-cy=create-test-button]').click();

  const $form = cy.get('[data-cy=create-test-modal] form');
  $form.should('be.visible');

  return $form;
};

export const deleteTest = () => {
  cy.location().then(({pathname}) => {
    const testId = pathname.split('/').pop();
    cy.visit('http://localhost:3000/');

    cy.get(`[data-cy=test-actions-button-${testId}]`).click();
    cy.get('[data-cy=test-delete-button]').click();

    cy.get(`[data-cy=test-actions-button-${testId}]`).should('not.exist');
  });
};
