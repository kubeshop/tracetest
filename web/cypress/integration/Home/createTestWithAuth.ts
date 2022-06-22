import {deleteTest} from '../utils/Common';

export function createTestWithAuth(
  $form: Cypress.Chainable<JQuery<HTMLElement>>,
  method: string,
  keys: string[],
  callback?: () => void
) {
  const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
  
  $form.get('[data-cy=create-test-name-input').type(name);
  $form.get('[data-cy=create-test-description-input').type(name);

  $form.get('[data-cy=create-test-next-button]').last().click();

  $form.get('[data-cy=url]').type('http://demo-pokemon-api.demo.svc.cluster.local/pokemon');
  $form.get('[data-cy=method-select]').click();
  $form.get('[data-cy=method-select-option-GET]').click();

  $form.get('[data-cy=auth-type-select]').click();
  $form.get(`[data-cy=auth-type-select-option-${method}]`).click();

  keys.forEach(key => {
    $form.get(`[data-cy=${method}-${key}]`).type(key);
  });

  if (callback) callback();

  $form.get('[data-cy=create-test-create-button]').last().click();

  cy.location('pathname').should('match', /\/test\/.*/i);
  cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
  deleteTest();
  return name;
}
