export const getAttributeListId = (number: number) => `#assertion-form_assertionList_${number}_attribute_list`;
export const getComparatorListId = (number: number) => `#assertion-form_assertionList_${number}_comparator_list`;

export function createAssertion(index = 0, brandNew = true) {
  cy.get(`[data-cy=trace-node-database]`, {timeout: 20000}).first().click({force: true});
  if (brandNew) {
    cy.get('[data-cy=add-assertion-button]').click({force: true});
  } else {
    cy.get('[data-cy=add-assertion-form-add-check]').click({force: true});
  }
  cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
  cy.get('[data-cy=assertion-check-attribute]').type('db');
  cy.wait(3500);
  const attributeListId = getAttributeListId(index);
  cy.get(`${attributeListId} + div .ant-select-item`).first().click({force: true});
  cy.wait(3500);
  cy.get('[data-cy=assertion-check-operator]').click({force: true});

  // const comparatorListId = getComparatorListId(index);
  // cy.get(`${comparatorListId} + div .ant-select-item`).last().click({force: true});

  // cy.get('[data-cy=assertion-check-value]').click({force: true});

  // cy.get('[data-cy=assertion-check-operator] + div .ant-select-selection-item').should('have.text', 'Contains');
  cy.get('[data-cy=assertion-form-submit-button]').click();
  cy.get('[data-cy=assertion-card-list]').should('be.visible');
  cy.get('[data-cy=assertion-card]').should('have.lengthOf', 1);
}
