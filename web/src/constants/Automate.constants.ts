export function CypressCodeSnippet(testName: string) {
  return `import Tracetest from '@tracetest/cypress';

const TRACETEST_API_TOKEN = Cypress.env('TRACETEST_API_TOKEN') || '';
const tracetest = Tracetest();

describe('Cypress Test', () => {
  before(done => {
    tracetest.configure(TRACETEST_API_TOKEN).then(() => done());
  });

  beforeEach(() => {
    cy.visit('/', {
      onBeforeLoad: win => tracetest.capture(win.document),
    });
  });

  afterEach(done => {
    tracetest.runTest('').then(() => done());
  });

  it('${testName}', () => {
    // ...cy commands
  });
});`;
}
