import {getAttributeListId, getValueFromList} from '../../support/commands';

describe('Test Run Detail Views', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('Trace view -> show the attribute list for a specific span', () => {
    cy.selectRunDetailMode(2);
    cy.get('[data-cy=trace-node-http]').click();
    // cy.get('[data-cy=toggle-drawer-SPAN_DETAILS]').click();

    cy.get('[data-cy=attribute-list]').should('be.visible');
    cy.get('[data-cy=attribute-row-http-method]').should('be.visible');
  });

  it('Trace view -> attribute list', () => {
    cy.selectRunDetailMode(2);
    cy.get('[data-cy=trace-node-http]').click();
    cy.get('[data-cy=toggle-drawer-SPAN_DETAILS]').click();

    cy.get('[data-cy=attribute-list]').should('be.visible');
  });

  it('Test view -> create assertion with empty selector', () => {
    cy.selectRunDetailMode(3);
    cy.get('[data-cy=add-test-spec-button]').click({force: true});

    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=editor-fallback]').should('not.exist');

    cy.get('[data-cy=selector-editor] [contenteditable]').first().clear();
    cy.get('[data-cy=selector-editor] .cm-placeholder').should('be.visible');

    cy.get('[data-cy=expression-editor] [contenteditable]').first().type('db.name', {delay: 100});
    cy.get(getAttributeListId(0)).first().click({force: true});

    // eslint-disable-next-line cypress/no-unnecessary-waiting
    cy.wait(100);
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).first().click();

    cy.get('[data-cy=assertion-form-submit-button]').click();
    cy.wait('@testRuns', {timeout: 30000});

    cy.intercept({method: 'PUT', url: '/api/tests/**'}).as('testPublish');
    cy.get('[data-cy=trace-actions-publish').click({force: true});
    cy.wait('@testPublish');
    cy.get('[data-cy=test-spec-container]', {timeout: 10000}).should('have.lengthOf', 1);
  });

  it('Test view -> create assertion with a nonexistent attribute (TDD).', () => {
    cy.selectRunDetailMode(3);
    cy.get('[data-cy=add-test-spec-button]').click({force: true});

    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=editor-fallback]').should('not.exist');

    cy.get('[data-cy=expression-editor] [contenteditable]').first().type('tdd', {delay: 100});
    cy.get('[data-cy=expression-editor] [contenteditable="true"]').last().type('value', {delay: 100});

    cy.get('[data-cy=assertion-form-submit-button]').click();
    cy.wait('@testRuns', {timeout: 30000});

    cy.intercept({method: 'PUT', url: '/api/tests/**'}).as('testPublish');
    cy.get('[data-cy=trace-actions-publish').click({force: true});
    cy.wait('@testPublish');
    cy.get('[data-cy=test-spec-container]', {timeout: 10000}).should('have.lengthOf', 1);
  });

  it('Test view -> navigate away with pending changes', () => {
    cy.selectRunDetailMode(3);

    cy.get('[data-cy=add-test-spec-button]').click({force: true});
    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=editor-fallback]').should('not.exist');

    cy.get('[data-cy=expression-editor] [contenteditable]').first().type('name', {delay: 100});
    cy.get(getAttributeListId(0)).first().click({force: true});

    // eslint-disable-next-line cypress/no-unnecessary-waiting
    cy.wait(100);
    cy.get('[data-cy=assertion-check-value] .cm-content').last().click();
    cy.get(getValueFromList(1)).first().click();

    cy.get('[data-cy=assertion-form-submit-button]').click();
    cy.wait('@testRuns', {timeout: 30000});
  });
});
