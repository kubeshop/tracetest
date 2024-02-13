const {default: Tracetest} = require('@tracetest/cypress');

const TRACETEST_API_TOKEN = Cypress.env('TRACETEST_API_TOKEN') || '';

let tracetest = undefined;

const definition = `
type: Test
spec:
  id: aW1wb3J0cyBhIHBva2Vtb24=
  name: imports a pokemon
  trigger:
    type: cypress
  specs:
  - selector: span[tracetest.span.type="http"] span[tracetest.span.type="http"]
    name: "All HTTP Spans: Status  code is 200"
    assertions:
    - attr:http.status_code   =   200
  - selector: span[tracetest.span.type="database"]
    name: "All Database Spans: Processing time is less than 100ms"
    assertions:
    - attr:tracetest.span.duration < 2s
  outputs:
  - name: MY_OUTPUT
    selector: span[tracetest.span.type="general" name="Tracetest trigger"]
    value: attr:name
`;

console.log('@@@Tracetest', Tracetest);

describe('Home', { defaultCommandTimeout: 80000 }, () => {
  before((done) => {
    Tracetest({ apiToken: TRACETEST_API_TOKEN }).then((instance) => {
      tracetest = instance;
      tracetest
        .setOptions({
          'imports a pokemon': {
            definition,
          },
        })
        .then(() => done());
    });
  });

  beforeEach(() => {
    cy.visit('/', {
      onBeforeLoad: (win) => tracetest.capture(win.document),
    });
  });

  // uncomment to wait for trace tests to be done
  after((done) => {
    tracetest.summary().then(() => done());
  });

  it('create a pokemon', () => {
    cy.get('[data-cy="create-pokemon-button"]').should('be.visible').click();
    cy.get('[data-cy="create-pokemon-modal"]').should('be.visible');
    cy.get('#name').type('Pikachu');
    cy.get('#type').type('Electric');
    cy.get('#imageUrl').type('https://oyster.ignimgs.com/mediawiki/apis.ign.com/pokemon-blue-version/8/89/Pikachu.jpg');

    cy.get('button').contains('OK').click();
  });

  it('imports a pokemon', () => {
    cy.get('[data-cy="import-pokemon-button"]').click();
    cy.get('[data-cy="import-pokemon-form"]').should('be.visible');

    cy.get('[id="id"]')
      .last()
      .type(Math.floor(Math.random() * 101).toString());
    cy.get('button').contains('OK').click({ force: true });
  });

  it('deletes a pokemon', () => {
    cy.get('[data-cy="pokemon-list"]').should('be.visible');
    cy.get('[data-cy="pokemon-card"]').first().click().get('[data-cy="delete-pokemon-button"]').first().click();
  });
});
