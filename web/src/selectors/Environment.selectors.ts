import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {endpoints} from 'redux/apis/TraceTest.api';
import UserSelectors from './User.selectors';

const stateSelector = (state: RootState) => state;

const selectEnvironmentList = createSelector(stateSelector, state => {
  const {data: {items = []} = {}} = endpoints.getEnvironments.select({})(state);

  return items;
});

const selectSelectedEnvironment = createSelector(selectEnvironmentList, stateSelector, (environmentList, state) => {
  const environmentId = UserSelectors.selectUserPreference(state, 'environmentId');

  return environmentList.find(({id}) => id === environmentId);
});

const EnvironmentSelectors = () => ({
  selectEnvironmentList,
  selectSelectedEnvironment,
  selectSelectedEnvironmentValues: createSelector(selectSelectedEnvironment, environment => {
    return environment?.values ?? [];
  }),
});

export default EnvironmentSelectors();
