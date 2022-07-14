export function fillCreateFormBasicStep(
  $form: Cypress.Chainable<JQuery<HTMLElement>>,
  name: string,
  description?: string
) {
  $form.get('[data-cy=create-test-next-button]').click();

  $form.get('[data-cy=create-test-name-input').type(name);
  $form.get('[data-cy=create-test-description-input').type(description || name);

  $form.get('[data-cy=create-test-next-button]').last().click();
}
