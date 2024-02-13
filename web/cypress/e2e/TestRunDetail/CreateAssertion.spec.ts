import {getAttributeListId, getValueFromList} from '../../support/commands';

describe('Create Assertion', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should create a basic assertion', () => {
    cy.createAssertion();
  });

  it('should create an assertion with multiple checks', () => {
    cy.selectRunDetailMode(3);

    cy.get(`[data-cy=trace-node-http]`, {timeout: 20000}).first().click();

    cy.get('[data-cy=add-test-spec-button]').click();
    cy.get('[data-cy=assertion-form]').should('be.visible');
    cy.get('[data-cy=editor-fallback]').should('not.exist');

    cy.get('[data-cy=expression-editor] [contenteditable]').first().type('http.status_code', {delay: 100});

    const attributeListId = getAttributeListId(0);
    cy.get(attributeListId).click();
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).last().click();

    cy.selectOperator(0);

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute] [contenteditable="true"]')
      .last()
      .type('tracetest.span.ty', {delay: 100});
    cy.get(getAttributeListId(0)).click();
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).last().click();

    cy.selectOperator(1);

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute] [contenteditable="true"]').last().type('duration', {delay: 100});
    cy.get(getAttributeListId(0)).click();
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).last().click();

    cy.selectOperator(2);

    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=test-spec-container]').should('have.lengthOf', 1);
  });

  it('should create a basic assertion using the advanced mode', () => {
    cy.selectRunDetailMode(3);

    cy.get(`[data-cy=trace-node-database]`, {timeout: 20000}).last().click();
    cy.get('[data-cy=add-test-spec-button]').click();
    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=editor-fallback]').should('not.exist');

    cy.get('[data-cy=editor-fallback]').should('not.exist');
    cy.get('[data-cy=selector-editor] [contenteditable]')
      .clear()
      .type('span[tracetest.span.type = "http"] span[tracetest.span.type = "database"]:first');

    cy.get('[data-cy=assertion-check-attribute] [contenteditable="true"]').type('db.name', {delay: 100});
    cy.get(getAttributeListId(0)).click();
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).first().click();

    cy.selectOperator(0);

    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=test-spec-container]').should('have.lengthOf', 1);
  });

  it('should update an assertion', () => {
    cy.createAssertion();
    cy.get(`[data-cy=edit-test-spec-button]`, {timeout: 20000});

    cy.get('[data-cy=edit-test-spec-button]').first().click();
    cy.get('[data-cy=assertion-form]').should('be.visible');
    cy.get('[data-cy=editor-fallback]').should('not.exist');

    cy.selectOperator(0);

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute] [contenteditable="true"]').last().type('duration', {delay: 100});
    cy.get(getAttributeListId(0)).click();
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).first().click();

    cy.selectOperator(1);

    cy.get('[data-cy=assertion-form-submit-button]').click();
    cy.get('[data-cy=test-spec-container]').should('have.lengthOf', 1);
  });

  it('should update an assertion with advanced mode', () => {
    cy.createAssertion();

    cy.get('[data-cy=edit-test-spec-button]').last().click();
    cy.get('[data-cy=assertion-form]').should('be.visible');
    cy.get('[data-cy=editor-fallback]').should('not.exist');

    cy.get('[data-cy=selector-editor] [contenteditable]').clear().type('span[tracetest.span.type = "database"]:last');

    cy.selectOperator(0);

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute] [contenteditable="true"]').last().type('duration', {delay: 100});
    cy.get(getAttributeListId(0)).click();
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).first().click();

    cy.selectOperator(1);

    cy.get('[data-cy=assertion-form-submit-button]').click();
    cy.get('[data-cy=test-spec-container]').should('have.lengthOf', 1);
  });

  it('should publish the changes', () => {
    cy.createAssertion();
    cy.get('[data-cy=trace-actions-publish').click({force: true});
    cy.get('[data-cy=test-spec-container]', {timeout: 10000}).should('have.lengthOf', 1);
  });

  it('should create an assertion and revert all changes', () => {
    cy.createAssertion();
    cy.get(`[data-cy=trace-node-database]`, {timeout: 20000}).last().click();
    cy.get('[data-cy=add-test-spec-button]').click();
    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=editor-fallback]').should('not.exist');

    cy.get('[data-cy=assertion-check-attribute] [contenteditable="true"]').type('db.name', {delay: 100});
    cy.get(getAttributeListId(0)).click();
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).first().click();

    cy.selectOperator(0);

    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=test-spec-container]').should('have.lengthOf', 2);
    cy.get('[data-cy=trace-actions-revert-all').click();
    cy.get('[data-cy=test-spec-container]').should('have.lengthOf', 0);
  });
});
