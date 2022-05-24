import {createAsyncThunk} from '@reduxjs/toolkit';
import TestDefinitionGateway from '../../gateways/TestDefinition.gateway';
import TestSelectors from '../../selectors/Test.selectors';
import TestDefinitionService from '../../services/TestDefinition.service';
import {TTestDefinitionEntry} from '../../types/TestDefinition.types';
import {RootState} from '../store';

const TestDefinitionActions = () => ({
  add: createAsyncThunk<void, {testId: string; definition: TTestDefinitionEntry}>(
    'testDefinition/add',
    async ({testId, definition}, {dispatch, getState}) => {
      const test = TestSelectors.selectTest(getState() as RootState, testId);

      const rawDefinitionList = test.definition.definitionList.map(def => TestDefinitionService.toRaw(def));

      await dispatch(
        TestDefinitionGateway.set(testId, {
          definitions: [...rawDefinitionList, TestDefinitionService.toRaw(definition)],
        })
      );
    }
  ),
  update: createAsyncThunk<void, {testId: string; selector: string; definition: TTestDefinitionEntry}>(
    'testDefinition/update',
    async ({testId, definition, selector}, {dispatch, getState}) => {
      const test = TestSelectors.selectTest(getState() as RootState, testId);
      const rawDefinitionList = test.definition.definitionList.map(def => TestDefinitionService.toRaw(def));

      await dispatch(
        TestDefinitionGateway.set(testId, {
          definitions: rawDefinitionList.map(def => {
            if (def.selector === selector) return TestDefinitionService.toRaw(definition);

            return def;
          }),
        })
      );
    }
  ),
  remove: createAsyncThunk<void, {testId: string; selector: string}>(
    'testDefinition/remove',
    async ({testId, selector}, {dispatch, getState}) => {
      const test = TestSelectors.selectTest(getState() as RootState, testId);
      const rawDefinitionList = test.definition.definitionList.map(def => TestDefinitionService.toRaw(def));

      await dispatch(
        TestDefinitionGateway.set(testId, {
          definitions: rawDefinitionList.filter(definition => definition.selector !== selector),
        })
      );
    }
  ),
});

export default TestDefinitionActions();
