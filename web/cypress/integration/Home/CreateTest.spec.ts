import {camelCase} from 'lodash';
import {DemoTestExampleList} from '../../../src/constants/Test.constants';

Cypress.on('uncaught:exception', err => !err.message.includes('ResizeObserver loop limit exceeded'));

const deleteTest = () => {
  cy.location().then(({pathname}) => {
    const testId = pathname.split('/').pop();
    cy.visit('http://localhost:3000/');

    cy.get(`[data-cy=test-actions-button-${testId}]`).click();
    cy.get('[data-cy=test-delete-button]').click();

    cy.get(`[data-cy=test-actions-button-${testId}]`).should('not.exist');
  });
};

const openCreateTestModal = () => {
  cy.get('[data-cy=create-test-button]').click();

  const $form = cy.get('[data-cy=create-test-modal] form');
  $form.should('be.visible');

  return $form;
};

beforeEach(() => {
  cy.visit('http://localhost:3000/');
});

it('should create a basic GET test from scratch', () => {
  const $form = openCreateTestModal();

  $form.get('[data-cy=method-select]').click();
  $form.get('[data-cy=method-select-option-GET]').click();
  const name = `Test - Shop - #${String(Date.now()).slice(-4)}`;

  $form.get('[data-cy=url]').type('http://shop/buy');
  $form.get('[data-cy=name').type(name);

  $form.get('[data-cy=create-test-submit]').click();

  cy.location('pathname').should('match', /\/test\/.*/i);
  cy.get('[data-cy=test-details-name]').should('have.text', name);
  deleteTest();
});

it('should create a basic POST test from scratch', () => {
  const $form = openCreateTestModal();

  $form.get('[data-cy=method-select]').click();
  $form.get('[data-cy=method-select-option-POST]').click();
  const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;

  $form.get('[data-cy=url]').type('http://demo-pokemon-api.demo.svc.cluster.local/pokemon');
  $form.get('[data-cy=name').type(name);
  $form
    .get('[data-cy=body]')
    .type(
      '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
      {
        parseSpecialCharSequences: false,
      }
    );

  $form.get('[data-cy=create-test-submit]').click();

  cy.location('pathname').should('match', /\/test\/.*/i);
  cy.get('[data-cy=test-details-name]').should('have.text', name);
  deleteTest();
});

it('should create a GET test from an example', () => {
  const [{name, description}] = DemoTestExampleList;

  const $form = openCreateTestModal();
  $form.get('[data-cy=example-button]').click();
  $form.get(`[data-cy=demo-example-${camelCase(name)}]`).click();
  $form.get('[data-cy=create-test-submit]').click();

  cy.location('pathname').should('match', /\/test\/.*/i);
  cy.get('[data-cy=test-details-name]').should('have.text', `${name} - ${description}`);
  deleteTest();
});

it('should create a POST test from an example', () => {
  const [, , {name, description}] = DemoTestExampleList;

  const $form = openCreateTestModal();
  $form.get('[data-cy=example-button]').click();
  $form.get(`[data-cy=demo-example-${camelCase(name)}]`).click();
  $form.get('[data-cy=create-test-submit]').click();

  cy.location('pathname').should('match', /\/test\/.*/i);
  cy.get('[data-cy=test-details-name]').should('have.text', `${name} - ${description}`);
  deleteTest();
});
