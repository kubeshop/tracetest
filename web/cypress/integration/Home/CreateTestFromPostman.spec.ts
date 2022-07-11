import {navigateToTestCreationPage} from '../utils/Common';
import {fillCreateFormBasicStep} from './fillCreateFormBasicStep';

describe('Create test from Postman Collection', () => {
  beforeEach(() => cy.visit('http://localhost:3000/'));

  it('should create a basic GET test', () => {
    const $form = navigateToTestCreationPage();

    $form.get('[data-cy=postman-plugin]').click();

    fillCreateFormBasicStep(
      $form,
      `Test - Pokemon - #${String(Date.now()).slice(-4)}`,
      'Create from Postman Collection'
    );

    cy.get('[data-cy="collectionFile"]').attachFile('collection.json');

    // sdkjfnds
  });
});
