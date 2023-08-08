Cypress.on('uncaught:exception', err => !err.message.includes('ResizeObserver loop limit exceeded'));

const testIdRegex = /\/test\/([\w-]+)/;
const testSuiteIdRegex = /\/testsuite\/([\w-]+)/;
const runIdRegex = /\/run\/(\w+)/;

export const getTestId = (pathname: string) => {
  cy.log(pathname);
  const result = pathname.match(testIdRegex);
  const testId = result.length > 1 ? result[1] : '';
  cy.log(testId);
  return testId;
};

export const getTestSuiteId = (pathname: string) => {
  cy.log(pathname);
  const result = pathname.match(testSuiteIdRegex);
  const testSuiteId = result.length > 1 ? result[1] : '';
  cy.log(testSuiteId);
  return testSuiteId;
};

export const getResultId = (pathname: string) => {
  cy.log(pathname);
  const result = pathname.match(runIdRegex);
  const runId = result.length > 1 ? result[1] : '';
  cy.log(runId);
  return runId;
};
