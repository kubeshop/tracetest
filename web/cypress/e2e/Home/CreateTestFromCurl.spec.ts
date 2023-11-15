import {POKEMON_HTTP_ENDPOINT} from '../constants/Test';

describe('Create test from CURL Command', () => {
  beforeEach(() => {
    cy.enableDemo();
    cy.visit('/');
  });

  it('should create a basic GET test', () => {
    cy.interceptHomeApiCall();
    cy.get('[data-cy=import-button]').click();
    cy.get('[data-cy=curl-plugin]').click();
    cy.get('[data-cy=curl-command-editor] [contenteditable]')
      .first()
      .type(
        `curl -XPOST ${POKEMON_HTTP_ENDPOINT}/pokemon
    -H Content-type: application/json
    --data '{"name":"meowth","type":"normal","imageUrl":"https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured":true}'
   `,
        {parseSpecialCharSequences: false}
      );
    cy.get(`[data-cy="import-test-submit"]`).click();
    cy.submitCreateForm();
    cy.matchTestRunPageUrl();
    cy.cancelOnBoarding();
    cy.get('[data-cy=overlay-input-overlay]').should('contain.text', POKEMON_HTTP_ENDPOINT);
    cy.deleteTest(true);
  });
});
