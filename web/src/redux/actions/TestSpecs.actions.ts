import {createAsyncThunk} from '@reduxjs/toolkit';
import {PatchCollection} from '@reduxjs/toolkit/dist/query/core/buildThunks';

import TestGateway from 'gateways/Test.gateway';
import TestRunGateway from 'gateways/TestRun.gateway';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import {selectTestOutputs} from 'redux/testOutputs/selectors';
import TestDefinitionService from 'services/TestDefinition.service';
import TestService from 'services/Test.service';
import TestRun from 'models/TestRun.model';
import Test from 'models/Test.model';
import AssertionResults from 'models/AssertionResults.model';
import {TTestSpecEntry} from 'models/TestSpecs.model';
import {RootState} from '../store';

export type TChange = {
  selector: string;
  action: 'add' | 'remove' | 'update';
  patch: PatchCollection;
};

const TestSpecsActions = () => ({
  publish: createAsyncThunk<TestRun, {test: Test; testId: string; runId: number}>(
    'testDefinition/publish',
    async ({test, testId, runId}, {dispatch, getState}) => {
      const specs = TestSpecsSelectors.selectSpecs(getState() as RootState).filter(def => !def.isDeleted);
      const outputs = selectTestOutputs(getState() as RootState);
      const rawTest = await TestService.getUpdatedRawTest(test, {definition: {specs}, outputs});
      await dispatch(TestGateway.edit(rawTest, testId));
      return dispatch(TestRunGateway.reRun(testId, runId)).unwrap();
    }
  ),
  dryRun: createAsyncThunk<AssertionResults, {definitionList: TTestSpecEntry[]; testId: string; runId: number}>(
    'testDefinition/dryRun',
    ({definitionList, testId, runId}, {dispatch}) => {
      const specs = definitionList.map(def => TestDefinitionService.toRaw(def));

      return dispatch(TestRunGateway.dryRun(testId, runId, {specs})).unwrap();
    }
  ),
});

export default TestSpecsActions();
