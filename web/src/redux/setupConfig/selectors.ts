import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';

const createTransactionSelectorsStateSelector = (state: RootState) => state.setupConfig;

const SetupConfigSelectors = () => ({
  selectDraftConfig: createSelector(createTransactionSelectorsStateSelector, ({draftConfig}) => draftConfig),
  selectIsFormValid: createSelector(createTransactionSelectorsStateSelector, ({isFormValid}) => isFormValid),
});

export default SetupConfigSelectors();
