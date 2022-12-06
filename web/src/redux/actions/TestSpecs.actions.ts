import {createAsyncThunk} from '@reduxjs/toolkit';
import {PatchCollection} from '@reduxjs/toolkit/dist/query/core/buildThunks';

import TestGateway from 'gateways/Test.gateway';
import TestRunGateway from 'gateways/TestRun.gateway';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import TestDefinitionService from 'services/TestDefinition.service';
import {TAssertionResults} from 'types/Assertion.types';
import {TTest} from 'types/Test.types';
import {TTestRun} from 'types/TestRun.types';
import {TTestSpecEntry} from 'types/TestSpecs.types';
import {RootState} from '../store';
import TestService from '../../services/Test.service';

export type TChange = {
  selector: string;
  action: 'add' | 'remove' | 'update';
  patch: PatchCollection;
};

const TestSpecsActions = () => ({
  publish: createAsyncThunk<TTestRun, {test: TTest; testId: string; runId: string}>(
    'testDefinition/publish',
    async ({test, testId, runId}, {dispatch, getState}) => {
      const specs = TestSpecsSelectors.selectSpecs(getState() as RootState).filter(def => !def.isDeleted);
      const rawTest = await TestService.getUpdatedRawTest(test, {definition: {specs}});
      await dispatch(TestGateway.edit(rawTest, testId));
      return dispatch(TestRunGateway.reRun(testId, runId)).unwrap();
    }
  ),
  dryRun: createAsyncThunk<TAssertionResults, {definitionList: TTestSpecEntry[]; testId: string; runId: string}>(
    'testDefinition/dryRun',
    ({definitionList, testId, runId}, {dispatch}) => {
      const specs = definitionList.map(def => TestDefinitionService.toRaw(def));

      return dispatch(TestRunGateway.dryRun(testId, runId, {specs})).unwrap();
    }
  ),
});

export default TestSpecsActions();
