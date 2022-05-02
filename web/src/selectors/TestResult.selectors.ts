import {createSelector} from '@reduxjs/toolkit';
import {RootState} from '../redux/store';

const stateSelector = (state: RootState) => state.resultList;

const TestResultSelectors = () => ({
  selectTestResultList(resultId = '') {
    return createSelector(stateSelector, ({resultListMap}) => resultListMap[resultId] || []);
  },
});

export default TestResultSelectors();
