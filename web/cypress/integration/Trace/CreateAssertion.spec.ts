import {createAssertion, getAttributeListId, getComparatorListId} from './CreateAssertion';

describe('Create Assertion', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should create a basic assertion', () => createAssertion());

  it('should create an assertion with multiple checks', () => {
    cy.get(`[data-cy=trace-node-http]`, {timeout: 20000}).first().click();

    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]').should('be.visible');

    cy.get('[data-cy=assertion-check-attribute]').type('http');
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
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 1);
  });

  it('should create a basic assertion using the advanced mode', () => {
    cy.get(`[data-cy=trace-node-database]`, {timeout: 20000}).last().click();
    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');

    cy.get('[data-cy=advanced-selector] [contenteditable]')
      .clear()
      .type('span[tracetest.span.type = "http"] span[tracetest.span.type = "database"]:first');

    cy.get('[data-cy=assertion-check-attribute]').type('db');
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
  //
  it('should update an assertion', () => {
    createAssertion();
    cy.get(`[data-cy=edit-assertion-button]`, {timeout: 20000});

    cy.get('[data-cy=edit-assertion-button]').first().click();
    cy.get('[data-cy=assertion-form]').should('be.visible');

    cy.get('[data-cy=assertion-check-operator]').first().click();
    cy.get(`${getComparatorListId(0)} + div .ant-select-item`)
      .first()
      .click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').first().should('have.text', 'Equals');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute]').last().type('service');
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
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 1);
  });

  it('should update an assertion with advanced mode', () => {
    createAssertion();

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
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 1);
  });
  it('should publish the changes', () => {
    createAssertion();
    cy.get('[data-cy=trace-actions-publish').click({force: true});
    cy.get('[data-cy=assertion-card]', {timeout: 10000}).should('have.lengthOf', 1);
  });

  it('should create an assertion and revert all changes', () => {
    createAssertion();
    cy.get(`[data-cy=trace-node-database]`, {timeout: 20000}).last().click();
    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');

    cy.get('[data-cy=assertion-check-attribute]').type('db');
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
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 2);

    cy.get('[data-cy=trace-actions-revert-all').click();
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 0);
  });
});
