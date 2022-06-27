import {openCreateTestModal} from '../utils/Common';
import {createTestWithAuth} from './createTestWithAuth';

describe('Create test with authentication', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/');
  });

  it('should create a basic GET test with api key authentication method', () => {
    const $form = openCreateTestModal();
    $form.get('[data-cy=create-test-next-button]').click();

    createTestWithAuth($form, 'apiKey', ['key', 'value'], () => {
      $form.get('[data-cy=auth-apiKey-select]').click();
      $form.get(`[data-cy=auth-apiKey-select-option-header]`).click();
    });
  });

  it('should create a basic GET test with basic authentication method', () => {
    const $form = openCreateTestModal();
    $form.get('[data-cy=create-test-next-button]').click();

    createTestWithAuth($form, 'basic', ['username', 'password']);
  });

  it('should create a basic GET test with bearer authentication method', () => {
    const $form = openCreateTestModal();
    $form.get('[data-cy=create-test-next-button]').click();

    createTestWithAuth($form, 'bearer', ['token']);
  });
});
