import {SupportedPlugins} from '../../../src/constants/Common.constants';
import {POKEMON_HTTP_ENDPOINT, PokeshopDemo} from '../constants/Test';

describe('Create test', () => {
  beforeEach(() => {
    cy.interceptHomeApiCall();
    cy.enableDemo();
    cy.visit('/');
  });
  afterEach(() => cy.deleteTest(true));

  it('should create a basic GET test from scratch', () => {
    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    cy.openTestCreationModal();
    cy.get(`[data-cy=${SupportedPlugins.REST.toLowerCase()}-plugin]`).click();
    cy.fillCreateFormBasicStep(name);

    cy.setCreateFormUrl('GET', `${POKEMON_HTTP_ENDPOINT}/pokemon`);
    cy.submitCreateForm();
    cy.makeSureUserIsOnTracePage();
    cy.get('[data-cy=overlay-input-overlay]').should('contain.text', name);
  });

  it('should create a basic POST test from scratch', () => {
    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    cy.openTestCreationModal();
    cy.get(`[data-cy=${SupportedPlugins.REST.toLowerCase()}-plugin]`).click();
    cy.fillCreateFormBasicStep(name);
    cy.setCreateFormUrl('POST', `${POKEMON_HTTP_ENDPOINT}/pokemon`);

    cy.get('#rc-tabs-0-tab-body').click();
    cy.get('[data-cy="bodyMode"]').click();
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
    cy.get('[data-cy=overlay-input-overlay]').should('contain.text', name);
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
