import {camelCase} from 'lodash';
import {Plugins} from '../../../src/constants/Plugins.constants';
import {deleteTest, navigateToTestCreationPage} from '../utils/Common';
import {fillCreateFormBasicStep} from './fillCreateFormBasicStep';

const DemoTestExampleList = Plugins.REST.demoList;

describe('Create test', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/');
  });

  it('should create a basic GET test from scratch', () => {
    const $form = navigateToTestCreationPage();
    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;

    fillCreateFormBasicStep($form, name);

    $form.get('[data-cy=method-select]').click();
    $form.get('[data-cy=method-select-option-GET]').click();
    $form.get('[data-cy=url]').type('http://demo-pokemon-api.demo.svc.cluster.local/pokemon');

    $form.get('[data-cy=create-test-create-button]').last().click();

    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
    deleteTest();
  });

  it('should create a basic POST test from scratch', () => {
    const $form = navigateToTestCreationPage();
    $form.get('[data-cy=create-test-next-button]').last().click();

    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;

    $form.get('[data-cy=create-test-name-input').type(name);
    $form.get('[data-cy=create-test-description-input').type(name);

    $form.get('[data-cy=create-test-next-button]').last().click();

    $form.get('[data-cy=method-select]').click();
    $form.get('[data-cy=method-select-option-POST]').click();
    $form.get('[data-cy=url]').type('http://demo-pokemon-api.demo.svc.cluster.local/pokemon');
    $form
      .get('[data-cy=body]')
      .type(
        '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
        {
          parseSpecialCharSequences: false,
        }
      );

    $form.get('[data-cy=create-test-create-button]').last().click();

    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
    deleteTest();
  });

  it('should create a GET test from an example', () => {
    const [{name}] = DemoTestExampleList;

    const $form = navigateToTestCreationPage();
    $form.get('[data-cy=create-test-next-button]').last().click();

    $form.get('[data-cy=example-button]').click();
    $form.get(`[data-cy=demo-example-${camelCase(name)}]`).click();

    $form.get('[data-cy=create-test-next-button]').last().click();
    $form.get('[data-cy=create-test-create-button]').last().click();

    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
    deleteTest();
  });

  it('should create a POST test from an example', () => {
    const [, {name}] = DemoTestExampleList;

    const $form = navigateToTestCreationPage();
    $form.get('[data-cy=create-test-next-button]').last().click();

    $form.get('[data-cy=example-button]').click();
    $form.get(`[data-cy=demo-example-${camelCase(name)}]`).click();

    $form.get('[data-cy=create-test-next-button]').last().click();
    $form.get('[data-cy=create-test-create-button]').last().click();

    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
    deleteTest();
  });
});
