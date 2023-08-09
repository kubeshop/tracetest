import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';

const createTestSuiteSelectorsStateSelector = (state: RootState) => state.createTestSuite;

const CreateTestSuiteSelectors = () => ({
  selectStepNumber: createSelector(createTestSuiteSelectorsStateSelector, ({stepNumber}) => stepNumber),
  selectDraft: createSelector(createTestSuiteSelectorsStateSelector, ({draft}) => draft),
  selectStepList: createSelector(createTestSuiteSelectorsStateSelector, ({stepList}) => stepList),
  selectIsFormValid: createSelector(createTestSuiteSelectorsStateSelector, ({isFormValid}) => isFormValid),
});

export default CreateTestSuiteSelectors();
