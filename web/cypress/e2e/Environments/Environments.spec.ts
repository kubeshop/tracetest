import {POKEMON_HTTP_ENDPOINT} from '../constants/Test';

interface IEnvironment {
  name: string;
  description: string;
  values: {key: string; value: string}[];
}

function createEnvironment(environment: IEnvironment) {
  cy.visit('/environments');
  cy.contains('Create Environment').click();

  cy.get('form#environment').within(() => {
    cy.get('#environment_name').type(environment.name);
    cy.get('#environment_description').type(environment.description);
    cy.get('[data-cy=interpolation-editor] [contenteditable]').eq(0).type(environment.values[0].key);
    cy.get('[data-cy=interpolation-editor] [contenteditable]').eq(1).type(environment.values[0].value);
  });

  cy.intercept({method: 'POST', url: '/api/environments'}).as('createEnvironment');
  cy.intercept({method: 'GET', url: '/api/environments?take=20&skip=0*'}).as('getEnvironments');

  cy.get('.ant-modal-footer').contains('Create').click();
  cy.wait('@createEnvironment');
  cy.wait('@getEnvironments');
}

function deleteEnvironment() {
  cy.visit('/environments');
  cy.get('[data-cy=environment-actions]').first().click();
  cy.get('[data-cy=environment-actions-delete]').click();

  cy.intercept({method: 'DELETE', url: '/api/environments/*'}).as('deleteEnvironment');
  cy.intercept({method: 'GET', url: '/api/environments?take=20&skip=0*'}).as('getEnvironments');

  cy.get('[data-cy=confirmation-modal] .ant-btn-primary').click();
  cy.wait('@deleteEnvironment');
  cy.wait('@getEnvironments');
}

describe('Environments', () => {
  const environment1: IEnvironment = {
    name: 'Environment One',
    description: 'Description Environment One',
    values: [{key: 'host', value: POKEMON_HTTP_ENDPOINT}],
  };
  const environment2: IEnvironment = {
    name: 'Environment Two',
    description: 'Description Environment Two',
    values: [{key: 'item2', value: 'value2'}],
  };

  it('should load the environments page', () => {
    cy.visit('/environments');
    cy.contains('All Environments');
    cy.get('input[placeholder*="Search environment"]');
    cy.contains('Create Environment');
  });

  it('should create a new environment', () => {
    createEnvironment(environment1);

    cy.contains(environment1.name).click();
    cy.contains(environment1.description);
    cy.contains(environment1.values[0].key);
    cy.contains(environment1.values[0].value);

    deleteEnvironment();
  });

  it('should update an environment', () => {
    createEnvironment(environment1);

    cy.get('[data-cy=environment-actions]').first().click();
    cy.get('[data-cy=environment-actions-edit]').click();

    cy.get('form#environment').within(() => {
      cy.get('#environment_name').clear().type(environment2.name);
      cy.get('#environment_description').clear().type(environment2.description);
      cy.get('[data-cy=interpolation-editor] [contenteditable]').eq(0).clear().type(environment2.values[0].key);
      cy.get('[data-cy=interpolation-editor] [contenteditable]').eq(1).clear().type(environment2.values[0].value);
    });

    cy.intercept({method: 'GET', url: '/api/environments?take=20&skip=0*'}).as('getEnvironments');

    cy.get('.ant-modal-footer').contains('Update').click();
    cy.wait('@getEnvironments');

    cy.contains(environment2.name).click();
    cy.contains(environment2.description);
    cy.contains(environment2.values[0].key);
    cy.contains(environment2.values[0].value);

    deleteEnvironment();
  });

  it('should delete an environment', () => {
    createEnvironment(environment1);
    deleteEnvironment();

    cy.contains(environment1.name).should('not.exist');
  });

  it('should create a test using variables from environment', () => {
    createEnvironment(environment1);

    cy.visit('/');
    cy.interceptHomeApiCall();

    // Select created environment
    cy.get('[data-cy=environment-selector]').click();
    cy.get('.environment-selector-items ul li').eq(1).click();

    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    cy.openTestCreationModal();
    cy.fillCreateFormBasicStep(name);
    // eslint-disable-next-line no-template-curly-in-string
    cy.setCreateFormUrl('GET', '${{}env:host}/pokemon');
    cy.submitCreateForm();
    cy.makeSureUserIsOnTracePage();

    cy.reload().get('[data-cy=run-detail-trigger-response]').within(() => {
      cy.contains('Environment').click();
      cy.contains(environment1.values[0].key);
      cy.contains(environment1.values[0].value);
    });

    cy.deleteTest(true);
    deleteEnvironment();
  });
});
