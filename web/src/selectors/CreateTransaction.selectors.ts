import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';

const createTransactionSelectorsStateSelector = (state: RootState) => state.createTransaction;

const CreateTestSelectors = () => ({
  selectStepNumber: createSelector(createTransactionSelectorsStateSelector, ({stepNumber}) => stepNumber),
  selectDraftTransaction: createSelector(createTransactionSelectorsStateSelector, ({draftTransaction}) => draftTransaction),
  selectStepList: createSelector(createTransactionSelectorsStateSelector, ({stepList}) => stepList),
  selectIsFormValid: createSelector(createTransactionSelectorsStateSelector, ({isFormValid}) => isFormValid),
});

export default CreateTestSelectors();
