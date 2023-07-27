import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {endpoints} from 'redux/apis/Tracetest';
import UserSelectors from './User.selectors';

const stateSelector = (state: RootState) => state;
const withOutputsSelector = (state: RootState, withOutputs: boolean) => withOutputs;
const selectOutputs = createSelector(stateSelector, ({testOutputs: {outputs}}) => {
  return outputs;
});

const selectEnvironmentList = createSelector(stateSelector, state => {
  const {data: {items = []} = {}} = endpoints.getEnvironments.select({})(state);

  return items;
});

const selectSelectedEnvironment = createSelector(selectEnvironmentList, stateSelector, (environmentList, state) => {
  const environmentId = UserSelectors.selectUserPreference(state, 'environmentId');

  return environmentList.find(({id}) => id === environmentId);
});

const selectSelectedEnvironmentValues = createSelector(selectSelectedEnvironment, environment => {
  return environment?.values ?? [];
});

const EnvironmentSelectors = () => ({
  selectEnvironmentList,
  selectSelectedEnvironment,
  withOutputsSelector,
  selectSelectedEnvironmentValues: createSelector(
    selectSelectedEnvironmentValues,
    selectOutputs,
    withOutputsSelector,
    (environmentValues, outputs, withOutputs) => {
      return environmentValues.concat(withOutputs ? outputs.map(({name, value}) => ({key: name, value})) : []);
    }
  ),
});

export default EnvironmentSelectors();
