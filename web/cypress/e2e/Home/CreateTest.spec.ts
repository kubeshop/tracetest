import {PokeshopDemo} from '../constants/Test';

describe('Create test', () => {
  beforeEach(() => {
    cy.interceptHomeApiCall();
    cy.visit('/');
  });
  afterEach(() => cy.deleteTest(true));

  it('should create a basic GET test from scratch', () => {
    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    cy.openTestCreationModal();
    cy.fillCreateFormBasicStep(name);
    cy.setCreateFormUrl('GET', 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon');
    cy.submitCreateForm();
    cy.makeSureUserIsOnTracePage();
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
  });

  it('should create a basic POST test from scratch', () => {
    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    cy.openTestCreationModal();
    cy.fillCreateFormBasicStep(name);
    cy.setCreateFormUrl('POST', 'http://demo-pokemon-api.demo.svc.cluster.local/pokemon');
    cy.get('[data-cy=bodyMode-json]').click();
    cy.get('[data-cy=body] [data-cy=interpolation-editor] [contenteditable]')
      .first()
      .type(
        '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}',
        {
          parseSpecialCharSequences: false,
        }
      );
    cy.submitCreateForm();
    cy.makeSureUserIsOnTracePage();
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
  });

  it('should create a GET test from an example', () => {
    const [{name}] = PokeshopDemo;
    cy.createTestByName(name);
  });

  it('should create a POST test from an example', () => {
    const [, {name}] = PokeshopDemo;
    cy.createTestByName(name);
  });
});
