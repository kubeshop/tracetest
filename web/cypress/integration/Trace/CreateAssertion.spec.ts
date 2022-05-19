import {createTest, deleteTest, testId} from '../utils/common';

describe('Create Assertion', () => {
  before(() => {
    createTest();
  });

  after(() => {
    deleteTest();
  });

  it('should create a basic assertion', () => {
    cy.visit(`http://localhost:3000/test/${testId}`);
    cy.get('[data-cy^=result-card]').first().click();

    cy.location('href').should('match', /\/test\/.*/i);

    cy.get(`[data-cy^=test-run-result-]`).first().click();
    cy.location('href').should('match', /\/result\/.*/i);

    cy.get('[data-cy^=trace-node-]', {timeout: 10000}).should('be.visible');
    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]', {timeout: 10000}).should('be.visible');

    // add selector
    cy.get('[data-cy=assertion-form-selector-input]').type('db');
    cy.get('.ant-select-item-option-content').first().click();
    cy.get('.ant-select-item-option-content').first().click();
    cy.get('.ant-select-item-option-content').first().click();

    cy.get('[data-cy=assertion-check-attribute]').type('db');
    cy.wait(500);
    cy.get('#assertion-form_assertionList_0_key_list + div .ant-select-item').first().click();

    cy.get('[data-cy=assertion-check-operator]').click();
    cy.get('#assertion-form_assertionList_0_compareOp_list + div .ant-select-item').last().click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').should('have.text', 'Contains');

    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('be.visible');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 1);
  });

  it('should create an assertion with multiple checks', () => {
    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=assertion-form]').should('be.visible');

    // add selector
    cy.get('[data-cy=assertion-form-selector-input]').type('db');
    cy.get('.ant-select-item-option-content').first().click();
    cy.get('.ant-select-item-option-content').first().click();
    cy.get('.ant-select-item-option-content').first().click();

    cy.get('[data-cy=assertion-form-selector-input]').type('service');
    cy.get('.ant-select-item-option-content').last().click();
    cy.get('.ant-select-item-option-content').last().click();
    cy.get('.ant-select-item-option-content').last().click();

    cy.get('[data-cy=assertion-check-attribute]').type('db');
    cy.wait(500);
    cy.get('#assertion-form_assertionList_0_key_list + div .ant-select-item').first().click();

    cy.get('[data-cy=assertion-check-operator]').click();
    cy.get('#assertion-form_assertionList_0_compareOp_list + div .ant-select-item').first().click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').first().should('have.text', 'Equals');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute]').last().type('service');
    cy.wait(500);
    cy.get('#assertion-form_assertionList_1_key_list + div .ant-select-item').first().click();

    cy.get('[data-cy=assertion-check-operator]').last().click();
    cy.get('#assertion-form_assertionList_1_compareOp_list + div .ant-select-item').last().click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').last().should('have.text', 'Contains');
    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('be.visible');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 2);
  });

  it('should update an assertion', () => {
    cy.get('[data-cy=edit-assertion-button]').last().click();
    cy.get('[data-cy=assertion-form]').should('be.visible');

    cy.get('[data-cy=assertion-check-operator]').first().click();
    cy.get('#assertion-form_assertionList_0_compareOp_list + div .ant-select-item:nth-child(2)').first().click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').first().should('have.text', 'Not equals');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-attribute]').last().type('service');
    cy.wait(500);
    cy.get('#assertion-form_assertionList_2_key_list + div .ant-select-item').first().click();

    cy.get('[data-cy=assertion-check-operator]').last().click();
    cy.get('#assertion-form_assertionList_2_compareOp_list + div .ant-select-item').last().click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').last().should('have.text', 'Contains');
    cy.get('[data-cy=assertion-form-submit-button]').click();

    cy.get('[data-cy=assertion-card-list]').should('be.visible');
    cy.get('[data-cy=assertion-card]').should('have.lengthOf', 2);
  });
});
