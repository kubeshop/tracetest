import {SupportedPlugins} from '../../../src/constants/Common.constants';
import {POKEMON_HTTP_ENDPOINT} from '../constants/Test';

describe('Create test from CURL Command', () => {
  beforeEach(() => {
    cy.enableDemo();
    cy.visit('/');
  });

  it('should create a basic GET test', () => {
    cy.interceptHomeApiCall();
    const name = `Test - Pokemon - #${String(Date.now()).slice(-4)}`;
    cy.openTestCreationModal();
    cy.get(`[data-cy=${SupportedPlugins.CURL.toLowerCase()}-plugin]`).click();
    cy.fillCreateFormBasicStep(name, 'Create from Curl Command');

    cy.get('[data-cy=curl-command-editor] [contenteditable]')
      .first()
      .type(
        `curl -XPOST ${POKEMON_HTTP_ENDPOINT}/pokemon
    -H Content-type: application/json
    --data '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}'
   `,
        {parseSpecialCharSequences: false}
      );
    cy.get('[data-cy=CreateTestFactory-create-next-button]').last().click();

    cy.submitCreateForm();
    cy.matchTestRunPageUrl();
    cy.cancelOnBoarding();
    cy.get('[data-cy=test-details-name]').should('have.text', `${name} (v1)`);
    cy.deleteTest(true);
  });
});
