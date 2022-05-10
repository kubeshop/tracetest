import {DOCUMENTATION_URL, GITHUB_URL} from '../../../src/constants/Common.constants';
import {createTest, deleteTest} from '../utils/common';

describe('Home', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/');
  });

  it('should render the layout', () => {
    cy.get('[data-cy=logo]').should('have.attr', 'src', '/static/media/Logo.da806c8b6e530d5c1a7ec68e4fe407fa.svg');
    cy.get('[data-cy=documentation-link]').should('have.attr', 'href', DOCUMENTATION_URL);
    cy.get('[data-cy=github-link]').should('have.attr', 'href', GITHUB_URL);
    cy.get('[data-cy=onboarding-link]').should('be.visible');
  });

  it('should render the list of tests', () => {
    cy.get('[data-cy=create-test-button]').should('be.visible');

    cy.get('[data-cy=test-list]').should('be.visible');
  });

  it('should go to the test details after click', () => {
    createTest();

    cy.visit('http://localhost:3000/');
    cy.get('[data-cy=test-card]').first().click();
    cy.location('href').should('match', /\/test\/.*/i);
  });

  it('should run a test from the home page', () => {
    cy.visit('http://localhost:3000/');
    cy.get('[data-cy=test-run-button]').first().click();
    cy.location('href').should('match', /\/test\/.*/i);

    deleteTest();
  });
});
