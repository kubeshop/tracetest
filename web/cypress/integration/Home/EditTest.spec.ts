import {getTestId} from '../utils/Common';

describe('Edit Test', () => {
  beforeEach(() => cy.createTest());
  afterEach(() => cy.deleteTest());

  it('should edit a test', () => {
    cy.location('pathname').then(pathname => {
      const testId = getTestId(pathname);
      cy.wait('@testList');
      cy.editTestByTestId(testId);
    });
  });
});
