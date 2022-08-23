Cypress.on('uncaught:exception', err => !err.message.includes('ResizeObserver loop limit exceeded'));

export const getTestId = (pathname: string) => {
  cy.log(pathname);
  const [, , localTestId] = pathname.split('/').reverse();

  cy.log(localTestId);
  return localTestId;
};
export const getResultId = (pathname: string) => {
  cy.log(pathname);
  const [resultId] = pathname.split('/').reverse();
  cy.log(resultId);
  return resultId;
};
