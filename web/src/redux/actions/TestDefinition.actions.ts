import {createAsyncThunk} from '@reduxjs/toolkit';
import {PatchCollection} from '@reduxjs/toolkit/dist/query/core/buildThunks';
import TestDefinitionGateway from '../../gateways/TestDefinition.gateway';
import TestRunGateway from '../../gateways/TestRun.gateway';
import TestSpecsSelectors from '../../selectors/TestSpecs.selectors';
import TestDefinitionService from '../../services/TestDefinition.service';
import {TAssertionResults} from '../../types/Assertion.types';
import {TRawTestSpecEntry, TTestSpecEntry} from '../../types/TestSpecs.types';
import {TTestRun} from '../../types/TestRun.types';
import {RootState} from '../store';

export type TChange = {
  selector: string;
  action: 'add' | 'remove' | 'update';
  patch: PatchCollection;
};

export type TCrudResponse = {
  definitionList: TTestSpecEntry[];
  change: TChange;
};

const TestDefinitionActions = () => ({
  publish: createAsyncThunk<TTestRun, {testId: string; runId: string}>(
    'testDefinition/publish',
    async ({testId, runId}, {dispatch, getState}) => {
      const rawDefinitionList = TestSpecsSelectors.selectSpecs(getState() as RootState).reduce<TRawTestSpecEntry[]>(
        (list, def) => (!def.isDeleted ? list.concat([TestDefinitionService.toRaw(def)]) : list),
        []
      );

      await dispatch(
        TestDefinitionGateway.set(testId, {
          specs: rawDefinitionList,
        })
      );

      return dispatch(TestRunGateway.reRun(testId, runId)).unwrap();
    }
  ),
  dryRun: createAsyncThunk<TAssertionResults, {definitionList: TTestSpecEntry[]; testId: string; runId: string}>(
    'testDefinition/dryRun',
    ({definitionList, testId, runId}, {dispatch}) => {
      const rawDefinitionList = definitionList.map(def => TestDefinitionService.toRaw(def));

      return dispatch(TestRunGateway.dryRun(testId, runId, {specs: rawDefinitionList})).unwrap();
    }
  ),
});

export default TestDefinitionActions();
