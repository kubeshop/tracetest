import {createTest, deleteTest} from '../utils/Common';

const getAttributeListId = (number: number) => `#assertion-form_assertionList_${number}_attribute_list`;
const getComparatorListId = (number: number) => `#assertion-form_assertionList_${number}_comparator_list`;

describe('Create Assertion', () => {
  let testId;
  before(() => {
    testId = createTest()
  });

  after(() => deleteTest());

  it('should create a basic assertion', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);
    cy.get('[data-cy^=result-card]', {timeout: 20000}).first().click();

    cy.location('href').should('match', /\/test\/.*/i);

    cy.get(`[data-cy^=test-run-result-]`).first().click();
    cy.location('href').should('match', /\/run\/.*/i);

    cy.get('[data-cy^=trace-node-]', {timeout: 20000}).should('be.visible');
    cy.get(`[data-cy=trace-node-database]`, {timeout: 20000}).first().click();

    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]', {timeout: 20000}).should('be.visible');

    cy.get('[data-cy=assertion-check-attribute]').type('db');
    cy.wait(500);
    cy.get(`${getAttributeListId(0)} + div .ant-select-item`)
      .first()
      .click();

    cy.get('[data-cy=assertion-check-operator]').click();
    cy.get(`${getComparatorListId(0)} + div .ant-select-item`)
      .last()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').should('have.text', 'Contains');

    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('be.visible');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 1);
  });

  it('should create an assertion with multiple checks', () => {
    cy.get(`[data-cy=trace-node-http]`).first().click();

    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]').should('be.visible');

    cy.get('[data-cy=assertion-check-attribute]').type('http');
    cy.wait(500);
    cy.get(`${getAttributeListId(0)} + div .ant-select-item`)
      .first()
      .click();

    cy.get('[data-cy=assertion-check-operator]').click();
    cy.get(`${getComparatorListId(0)} + div .ant-select-item`)
      .first()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').first().should('have.text', 'Equals');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute]').last().type('service');
    cy.wait(500);
    cy.get(`${getAttributeListId(1)} + div .ant-select-item`)
      .first()
      .click();

    cy.get('[data-cy=assertion-check-operator]').last().click();
    cy.get(`${getComparatorListId(1)} + div .ant-select-item`)
      .last()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').last().should('have.text', 'Contains');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute]').last().type('duration');
    cy.wait(500);
    cy.get(`${getAttributeListId(2)} + div .ant-select-item`)
      .first()
      .click();

    cy.get('[data-cy=assertion-check-operator]').last().click();
    cy.get(`${getComparatorListId(2)} + div .ant-select-item`)
      .last()
      .click();

    cy.get('[data-cy=assertion-check-value]').last().type('s');
    cy.get('[data-cy=duration]').click();
    cy.get(`[data-cy=duration-unit-μs]`).click();

    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('be.visible');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 2);
  });

  it('should create a basic assertion using the advanced mode', () => {
    cy.get(`[data-cy=trace-node-database]`).last().click();
    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]', {timeout: 20000}).should('be.visible');

    cy.get('[data-cy=mode-selector-switch]').click();
    cy.get('[data-cy=advanced-selector] [contenteditable]')
      .clear()
      .type('span[tracetest.span.type = "http"] span[tracetest.span.type = "database"]:first');

    cy.get('[data-cy=assertion-check-attribute]').type('db');
    cy.wait(500);
    cy.get(`${getAttributeListId(0)} + div .ant-select-item`)
      .first()
      .click();

    cy.get('[data-cy=assertion-check-operator]').click();
    cy.get(`${getComparatorListId(0)} + div .ant-select-item`)
      .last()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').should('have.text', 'Contains');

    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('be.visible');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 3);
  });

  it('should update an assertion', () => {
    cy.get('[data-cy=edit-assertion-button]').first().click();
    cy.get('[data-cy=assertion-form]').should('be.visible');

    cy.get('[data-cy=assertion-check-operator]').first().click();
    cy.get(`${getComparatorListId(0)} + div .ant-select-item`)
      .first()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').first().should('have.text', 'Equals');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute]').last().type('service');
    cy.wait(500);
    cy.get(`${getAttributeListId(1)} + div .ant-select-item`)
      .first()
      .click();

    cy.get('[data-cy=assertion-check-operator]').last().click();
    cy.get(`${getComparatorListId(1)} + div .ant-select-item`)
      .last()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').last().should('have.text', 'Contains');
    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('be.visible');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 3);
  });

  it('should update an assertion with advanced mode', () => {
    cy.get('[data-cy=edit-assertion-button]').last().click();
    cy.get('[data-cy=assertion-form]').should('be.visible');

    cy.get('[data-cy=advanced-selector] [contenteditable]').clear().type('span[tracetest.span.type = "database"]:last');

    cy.get('[data-cy=assertion-check-operator]').first().click();
    cy.get(`${getComparatorListId(0)} + div .ant-select-item`)
      .first()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').first().should('have.text', 'Equals');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute]').last().type('service');
    cy.wait(500);
    cy.get(`${getAttributeListId(1)} + div .ant-select-item`)
      .first()
      .click();

    cy.get('[data-cy=assertion-check-operator]').last().click();
    cy.get(`${getComparatorListId(1)} + div .ant-select-item`)
      .last()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').last().should('have.text', 'Contains');
    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('be.visible');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 3);
  });

  it('should publish the changes', () => {
    cy.get('[data-cy=trace-actions-publish').click();
    cy.get('[data-cy=assertion-card]', {timeout: 20000}).should('have.lengthOf', 3);
  });

  it('should create an assertion and revert all changes', () => {
    cy.get(`[data-cy=trace-node-database]`, {timeout: 20000}).last().click();
    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]', {timeout: 20000}).should('be.visible');

    cy.get('[data-cy=assertion-check-attribute]').type('db');
    cy.wait(500);
    cy.get(`${getAttributeListId(0)} + div .ant-select-item`)
      .first()
      .click();

    cy.get('[data-cy=assertion-check-operator]').click();
    cy.get(`${getComparatorListId(0)} + div .ant-select-item`)
      .last()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').should('have.text', 'Contains');

    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('exist');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 4);

    cy.get('[data-cy=trace-actions-revert-all').click();
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 3);
  });
});
