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

    cy.wait(7000);
    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=create-assertion-form]', {timeout: 10000}).should('be.visible');

    cy.get('[data-cy=item-selector-tag] + [role=img]').first().click();
    cy.get('[data-cy=affected-spans-count]')
      .invoke('text')
      .should('match', /Affects \d+ spans/);

    cy.get('[data-cy=assertion-check-key]').type('http');
    cy.get('.ant-select-item-option-content').first().click();
    cy.get('[data-cy=assertion-check-key]').should('have.text', 'http.status_code');
    cy.get('[data-cy=assertion-check-value] input').should('have.attr', 'value', '200');

    cy.get('[data-cy=assertion-check-operator]').click();
    cy.get('.ant-select-item-option-content').last().click();

    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').should('have.text', 'contains');
    cy.get('#add-assertion-modal-ok-button').click();

    cy.get('[data-cy=assertion-table]').should('be.visible');
  });

  it('should create an assertion with multiple checks', () => {
    cy.get('[data-cy=add-assertion-button]').click();
    cy.get('[data-cy=create-assertion-form]').should('be.visible');

    cy.get('[data-cy=assertion-check-key]').first().type('http');
    cy.get('.ant-select-item-option-content').first().click();
    cy.get('[data-cy=assertion-check-key]').first().should('have.text', 'http.status_code');
    cy.get('[data-cy=assertion-check-value] input').first().should('have.attr', 'value', '200');
    cy.get('[data-cy=assertion-check-operator]').first().click();
    cy.get('.ant-select-item-option-content').last().click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').first().should('have.text', 'contains');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-key]').last().type('service');
    cy.get('.ant-select-item-option-content').last().click();
    cy.get('[data-cy=assertion-check-key]').last().should('have.text', 'service.name');
    cy.get('[data-cy=assertion-check-value] input').last().should('have.attr', 'value', 'pokeshop');

    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').last().should('have.text', 'eq');
    cy.get('#add-assertion-modal-ok-button').click();

    cy.get('[data-cy=assertion-table]').should('have.lengthOf', 2);
  });

  it('should update an assertion', () => {
    cy.get('[data-cy=edit-assertion-button]').last().click();
    cy.get('[data-cy=create-assertion-form]').should('be.visible');

    cy.get('[data-cy=assertion-check-operator]').first().click();
    cy.get('.ant-select-item-option-content').first().click();
    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').first().should('have.text', 'eq');

    cy.get('[data-cy=add-assertion-form-add-check]').click();

    cy.get('[data-cy=assertion-check-key]').last().type('service');
    cy.get('.ant-select-item-option-content').last().click();
    cy.get('[data-cy=assertion-check-key]').last().should('have.text', 'service.name');
    cy.get('[data-cy=assertion-check-value] input').last().should('have.attr', 'value', 'pokeshop');

    cy.get('[data-cy=assertion-check-operator] .ant-select-selection-item').last().should('have.text', 'eq');
    cy.get('#add-assertion-modal-ok-button').click();

    cy.get('[data-cy=assertion-table]').should('have.lengthOf', 2);

    cy.get('[data-cy=test-results-assertion-table]').should('have.lengthOf', 2);
  });
});
