import {getTestId} from '../utils/Common';

describe('Edit Test', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should edit a test', () => {
    cy.location('pathname').then(pathname => {
      const testId = getTestId(pathname);
      cy.visit(`/`);
      cy.wait('@testList');
      cy.editTestByTestId(testId);
    });
  });

  it('should edit a test from the test details', () => {
    cy.location('pathname').then(pathname => {
      const testId = getTestId(pathname);
      cy.visit(`/`);
      cy.wait('@testList');
      cy.visit(`/test/${testId}`);
      cy.editTestByTestId(testId);
    });
  });
});
