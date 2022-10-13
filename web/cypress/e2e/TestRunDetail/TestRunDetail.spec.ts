import {getAttributeListId} from '../../support/commands';

describe('Test Run Detail Views', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('Trace view -> show the attribute list for a specific span', () => {
    cy.selectRunDetailMode(2);
    cy.get('[data-cy=trace-node-http]').click();

    cy.get('[data-cy=attribute-list]').should('be.visible');
    cy.get('[data-cy=attribute-row-http-method]').should('be.visible');
  });
  it('Trace view -> attribute list -> switch between tabs', () => {
    cy.selectRunDetailMode(2);
    cy.get('[data-cy=trace-node-http]').click();

    cy.get('[data-cy=attribute-list]').should('be.visible');
    cy.get('[data-cy=attribute-tabs-response]').should('be.visible').click();
  });
  it('Test view -> create assertion with empty selector', () => {
    cy.selectRunDetailMode(3);
    cy.get('[data-cy=add-test-spec-button]').click({force: true});

    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=advanced-selector] [contenteditable]').clear();
    cy.get('[data-cy=assertion-check-attribute]').type('db.name');
    const attributeListId = getAttributeListId(0);
    cy.get(`${attributeListId} + div .ant-select-item:nth-child(2)`).first().click({force: true});
    cy.get('[data-cy=assertion-form-submit-button]').click();
    cy.wait('@testRuns', {timeout: 30000});

    cy.intercept({method: 'PUT', url: '/api/tests/**/definition'}).as('testPublish');
    cy.get('[data-cy=trace-actions-publish').click({force: true});
    cy.wait('@testPublish');
    cy.get('[data-cy=test-spec-container]', {timeout: 10000}).should('have.lengthOf', 1);
  });
  it('Test view -> create assertion with a nonexistent attribute (TDD).', () => {
    cy.selectRunDetailMode(3);
    cy.get('[data-cy=add-test-spec-button]').click({force: true});

    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');

    cy.get('[data-cy=assertion-check-attribute]').type('tdd').click();
    const attributeListId = getAttributeListId(0);
    cy.get(`${attributeListId} + div .ant-select-item:nth-child(1)`).first().click({force: true});
    cy.get('[data-cy=assertion-check-value]').type('value');

    cy.get('[data-cy=assertion-form-submit-button]').click();
    cy.wait('@testRuns', {timeout: 30000});

    cy.intercept({method: 'PUT', url: '/api/tests/**/definition'}).as('testPublish');
    cy.get('[data-cy=trace-actions-publish').click({force: true});
    cy.wait('@testPublish');
    cy.get('[data-cy=test-spec-container]', {timeout: 10000}).should('have.lengthOf', 1);
  });
  it('Test view -> navigate away with pending changes', () => {
    cy.selectRunDetailMode(3);
    cy.get('[data-cy=add-test-spec-button]').click({force: true});
    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=assertion-check-attribute]').type('db').click();

    const attributeListId = getAttributeListId(0);
    cy.get(`${attributeListId} + div .ant-select-item:nth-child(1)`).first().click({force: true});

    cy.get('[data-cy=assertion-check-value]').type('value');
    cy.get('[data-cy=assertion-form-submit-button]').click();
    cy.wait('@testRuns', {timeout: 30000});
  });
});
