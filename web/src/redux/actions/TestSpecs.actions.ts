import {createAsyncThunk} from '@reduxjs/toolkit';
import {PatchCollection} from '@reduxjs/toolkit/dist/query/core/buildThunks';

import TestRunGateway from 'gateways/TestRun.gateway';
import TestSpecsGateway from 'gateways/TestSpecs.gateway';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import TestDefinitionService from 'services/TestDefinition.service';
import {TAssertionResults} from 'types/Assertion.types';
import {TTestRun} from 'types/TestRun.types';
import {TRawTestSpecEntry, TTestSpecEntry} from 'types/TestSpecs.types';
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

const TestSpecsActions = () => ({
  publish: createAsyncThunk<TTestRun, {testId: string; runId: string}>(
    'testDefinition/publish',
    async ({testId, runId}, {dispatch, getState}) => {
      const rawDefinitionList = TestSpecsSelectors.selectSpecs(getState() as RootState).reduce<TRawTestSpecEntry[]>(
        (list, def) => (!def.isDeleted ? list.concat([TestDefinitionService.toRaw(def)]) : list),
        []
      );
      const specs = TestDefinitionService.formatExpectedField(rawDefinitionList);

      await dispatch(TestSpecsGateway.set(testId, {specs}));

      return dispatch(TestRunGateway.reRun(testId, runId)).unwrap();
    }
  ),
  dryRun: createAsyncThunk<TAssertionResults, {definitionList: TTestSpecEntry[]; testId: string; runId: string}>(
    'testDefinition/dryRun',
    ({definitionList, testId, runId}, {dispatch}) => {
      const rawDefinitionList = definitionList.map(def => TestDefinitionService.toRaw(def));
      const specs = TestDefinitionService.formatExpectedField(rawDefinitionList);

      return dispatch(TestRunGateway.dryRun(testId, runId, {specs})).unwrap();
    }
  ),
});

export default TestSpecsActions();
