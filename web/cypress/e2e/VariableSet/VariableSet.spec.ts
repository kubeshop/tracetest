import {SupportedPlugins} from '../../../src/constants/Common.constants';
import {POKEMON_HTTP_ENDPOINT} from '../constants/Test';

interface IVariableSet {
  name: string;
  description: string;
  values: {key: string; value: string}[];
}

function createVariableSet(variableSet: IVariableSet) {
  cy.visit('/variablesets');
  cy.contains('Create Variable Set').click();

  cy.get('form#variableSet').within(() => {
    cy.get('#variableSet_name').type(variableSet.name);
    cy.get('#variableSet_description').type(variableSet.description);
    cy.get('[data-cy=interpolation-editor] [contenteditable]').eq(0).type(variableSet.values[0].key);
    cy.get('[data-cy=interpolation-editor] [contenteditable]').eq(1).type(variableSet.values[0].value);
  });

  cy.intercept({method: 'POST', url: '/api/variablesets'}).as('createVariableSet');
  cy.intercept({method: 'GET', url: '/api/variablesets?take=20&skip=0*'}).as('getVariableSet');

  cy.get('.ant-modal-footer').contains('Create').click();
  cy.wait('@createVariableSet');
  cy.wait('@getVariableSet');
}

function deleteVariableSet() {
  cy.visit('/variablesets');
  cy.get('[data-cy=variableSet-actions]').first().click();
  cy.get('[data-cy=variableSet-actions-delete]').click();

  cy.intercept({method: 'DELETE', url: '/api/variablesets/*'}).as('deleteVariableSet');
  cy.intercept({method: 'GET', url: '/api/variablesets?take=20&skip=0*'}).as('getVariableSet');

  cy.get('[data-cy=confirmation-modal] .ant-btn-primary').click();
  cy.wait('@deleteVariableSet');
  cy.wait('@getVariableSet');
}

describe('VariableSets', () => {
  const variableSet1: IVariableSet = {
    name: 'Vars One',
    description: 'Description variableSet One',
    values: [{key: 'host', value: POKEMON_HTTP_ENDPOINT}],
  };
  const variableSet2: IVariableSet = {
    name: 'Vars Two',
    description: 'Description variableSet Two',
    values: [{key: 'item2', value: 'value2'}],
  };

  it('should load the variableSets page', () => {
    cy.visit('/variablesets');
    cy.contains('All Variable Sets');
    cy.get('input[placeholder*="Search variable set"]');
    cy.contains('Create Variable Set');
  });

  it('should create a new variableSet', () => {
    createVariableSet(variableSet1);

    cy.contains(variableSet1.name).click();
    cy.contains(variableSet1.description);
    cy.contains(variableSet1.values[0].key);
    cy.contains(variableSet1.values[0].value);

    deleteVariableSet();
  });

  it('should update an variableSet', () => {
    createVariableSet(variableSet1);

    cy.get('[data-cy=variableSet-actions]').first().click();
    cy.get('[data-cy=variableSet-actions-edit]').click();

    cy.get('form#variableSet').within(() => {
      cy.get('#variableSet_name').clear().type(variableSet2.name);
      cy.get('#variableSet_description').clear().type(variableSet2.description);
      cy.get('[data-cy=interpolation-editor] [contenteditable]').eq(0).clear().type(variableSet2.values[0].key);
      cy.get('[data-cy=interpolation-editor] [contenteditable]').eq(1).clear().type(variableSet2.values[0].value);
    });

    cy.intercept({method: 'GET', url: '/api/variablesets?take=20&skip=0*'}).as('getVariableSets');

    cy.get('.ant-modal-footer').contains('Update').click();
    cy.wait('@getVariableSets');

    cy.contains(variableSet2.name).click();
    cy.contains(variableSet2.description);
    cy.contains(variableSet2.values[0].key);
    cy.contains(variableSet2.values[0].value);

    deleteVariableSet();
  });

  it('should delete an variableSet', () => {
    createVariableSet(variableSet1);
    deleteVariableSet();

    cy.contains(variableSet1.name).should('not.exist');
  });

  it.only('should create a test using variables from variableSet', () => {
    createVariableSet(variableSet1);

    cy.visit('/tests');
    cy.interceptHomeApiCall();

    // Select created variableSet
    cy.get('[data-cy=variableSet-selector]').click();
    cy.get('.variableSet-selector-items ul li').eq(1).click();

    cy.openTestCreationModal();
    cy.get(`[data-cy=${SupportedPlugins.REST.toLowerCase()}-plugin]`).click();
    // eslint-disable-next-line no-template-curly-in-string
    cy.setCreateFormUrl('GET', '${{}env:host}/pokemon');
    cy.submitCreateForm();
    cy.makeSureUserIsOnTracePage();

    cy.reload()
      .get('[data-cy=run-detail-trigger-response]')
      .within(() => {
        cy.contains('Variable Set').click();
        cy.contains(variableSet1.values[0].key);
        cy.contains(variableSet1.values[0].value);
      });

    cy.deleteTest(true);
    deleteVariableSet();
  });
});
