import {DOCUMENTATION_URL, GITHUB_URL} from '../../../src/constants/Common.constants';

Cypress.on('uncaught:exception', err => !err.message.includes('ResizeObserver loop limit exceeded'));

describe('Home', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/');
  });

  it('renders the layout', () => {
    cy.get('[data-cy=logo]').should('have.attr', 'src', '/static/media/Logo.e43c6126615ae2bbdf534a0338b1dc1d.svg');
    cy.get('[data-cy=documentation-link]').should('have.attr', 'href', DOCUMENTATION_URL);
    cy.get('[data-cy=github-link]').should('have.attr', 'href', GITHUB_URL);
    cy.get('[data-cy=onboarding-link]').should('be.visible');
  });

  it('renders the list of tests', () => {
    cy.get('[data-cy=create-test-button]').should('be.visible');

    cy.get('[data-cy=testList]').should('be.visible');
    cy.get('[data-cy=testList] .ant-table-row').should($tr => {
      expect($tr).to.have.length.greaterThan(0);
    });
  });
});
