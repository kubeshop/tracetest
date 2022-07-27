import {navigateToTestCreationPage} from '../utils/Common';
import {fillCreateFormBasicStep} from './fillCreateFormBasicStep';

describe('Create test from Postman Collection', () => {
  beforeEach(() => cy.visit('http://localhost:3000/'));

  it('should create a basic GET test', () => {
    const $form = navigateToTestCreationPage();

    $form.get('[data-cy=postman-plugin]').click();

    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    fillCreateFormBasicStep($form, name, 'Create from Postman Collection');

    cy.get('[data-cy="collectionFile"]').attachFile('collection.json');

    $form.get('[data-cy=collectionTest-select]').click();
    $form.get('[data-cy=collectionTest-1]').click({force: true});

    $form.get('[data-cy=create-test-create-button]').last().click();

    cy.location('pathname').should('match', /\/test\/.*/i);
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
  });
});
