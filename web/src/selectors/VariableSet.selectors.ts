import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import TracetestAPI from 'redux/apis/Tracetest';
import UserSelectors from './User.selectors';

const stateSelector = (state: RootState) => state;
const withOutputsSelector = (state: RootState, withOutputs: boolean) => withOutputs;
const selectOutputs = createSelector(stateSelector, ({testOutputs: {outputs}}) => {
  return outputs;
});

const selectVariableSetList = createSelector(stateSelector, state => {
  const {data: {items = []} = {}} = TracetestAPI.instance.endpoints.getVariableSets.select({})(state);

  return items;
});

const selectSelectedVariableSet = createSelector(selectVariableSetList, stateSelector, (variableSetList, state) => {
  const variableSetId = UserSelectors.selectUserPreference(state, 'variableSetId');

  return variableSetList.find(({id}) => id === variableSetId);
});

const selectSelectedVariableSetValues = createSelector(selectSelectedVariableSet, variableSet => {
  return variableSet?.values ?? [];
});

const VariableSetSelectors = () => ({
  selectVariableSetList,
  selectSelectedVariableSet,
  withOutputsSelector,
  selectSelectedVariableSetValues: createSelector(
    selectSelectedVariableSetValues,
    selectOutputs,
    withOutputsSelector,
    (variableSetValue, outputs, withOutputs) => {
      return variableSetValue.concat(withOutputs ? outputs.map(({name, value}) => ({key: name, value})) : []);
    }
  ),
});

export default VariableSetSelectors();
