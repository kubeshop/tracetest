import {createAsyncThunk} from '@reduxjs/toolkit';
import {PatchCollection} from '@reduxjs/toolkit/dist/query/core/buildThunks';
import TestDefinitionGateway from '../../gateways/TestDefinition.gateway';
import TestRunGateway from '../../gateways/TestRun.gateway';
import TestDefinitionSelectors from '../../selectors/TestDefinition.selectors';
import TestDefinitionService from '../../services/TestDefinition.service';
import {TAssertionResults} from '../../types/Assertion.types';
import {TRawTestDefinitionEntry, TTestDefinitionEntry} from '../../types/TestDefinition.types';
import {TTestRun} from '../../types/TestRun.types';
import {RootState} from '../store';

export type TChange = {
  selector: string;
  action: 'add' | 'remove' | 'update';
  patch: PatchCollection;
};

export type TCrudResponse = {
  definitionList: TTestDefinitionEntry[];
  change: TChange;
};

const TestDefinitionActions = () => ({
  publish: createAsyncThunk<TTestRun, {testId: string; runId: string}>(
    'testDefinition/publish',
    async ({testId, runId}, {dispatch, getState}) => {
      const rawDefinitionList = TestDefinitionSelectors.selectDefinitionList(getState() as RootState).reduce<
        TRawTestDefinitionEntry[]
      >((list, def) => (!def.isDeleted ? list.concat([TestDefinitionService.toRaw(def)]) : list), []);

      await dispatch(
        TestDefinitionGateway.set(testId, {
          definitions: rawDefinitionList,
        })
      );

      return dispatch(TestRunGateway.reRun(testId, runId)).unwrap();
    }
  ),
  dryRun: createAsyncThunk<TAssertionResults, {definitionList: TTestDefinitionEntry[]; testId: string; runId: string}>(
    'testDefinition/dryRun',
    ({definitionList, testId, runId}, {dispatch}) => {
      const rawDefinitionList = definitionList.map(def => TestDefinitionService.toRaw(def));

      return dispatch(TestRunGateway.dryRun(testId, runId, {definitions: rawDefinitionList})).unwrap();
    }
  ),
});

export default TestDefinitionActions();
