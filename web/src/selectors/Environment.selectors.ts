import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {endpoints} from 'redux/apis/TraceTest.api';
import UserSelectors from './User.selectors';

const stateSelector = (state: RootState) => state;

const selectEnvironmentList = createSelector(stateSelector, state => {
  const {data: {items = []} = {}} = endpoints.getEnvList.select({})(state);

  return items;
});

const EnvironmentSelectors = () => ({
  selectEnvironmentList,
  selectSelectedEnvironment: createSelector(selectEnvironmentList, stateSelector, (environmentList, state) => {
    const environmentId = UserSelectors.selectUserPreference(state, 'environmentId');

    return environmentList.find(({id}) => id === environmentId);
  }),
  selectSelectedEnvironmentEntryList: createSelector(selectEnvironmentList, stateSelector, (environmentList, state) => {
    const environmentId = UserSelectors.selectUserPreference(state, 'environmentId');
    const {data: entries = []} = endpoints.getEnvironmentSecretList.select({environmentId})(state);

    return entries;
  }),
});

export default EnvironmentSelectors();
