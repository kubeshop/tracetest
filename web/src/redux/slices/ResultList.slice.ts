import {createAsyncThunk, createSlice, PayloadAction} from '@reduxjs/toolkit';
import TestRunResultGateway from '../../gateways/TestRunResult.gateway';
import TraceService from '../../services/Trace.service';
import {IAssertionResult, TAssertionResultList} from '../../types/Assertion.types';
import {ITest} from '../../types/Test.types';
import {ITrace} from '../../types/Trace.types';

interface ITestResultListState {
  resultListMap: Record<string, IAssertionResult[]>;
}

export interface ITestResultListReplaceParams {
  resultList: IAssertionResult[];
  resultId: string;
}

const initialState: ITestResultListState = {
  resultListMap: {},
};

export const updateTestResult = createAsyncThunk<
  ITestResultListReplaceParams,
  {trace: ITrace; resultId: string; test: ITest}
>('resultList/load', async ({trace, test, resultId}, {dispatch}) => {
  const resultList = TraceService.runTest(trace, test);

  await dispatch(
    TestRunResultGateway.update(test.testId, resultId, TraceService.parseAssertionResultListToTestResult(resultList))
  );

  return {resultId, resultList};
});

const ResultListSlice = createSlice({
  name: 'resultList',
  initialState,
  reducers: {
    replace(
      state,
      {
        payload: {assertionResult, test, trace, resultId},
      }: PayloadAction<{assertionResult: TAssertionResultList; test: ITest; trace: ITrace; resultId: string}>
    ) {
      const resultList = TraceService.parseTestResultToAssertionResultList(assertionResult, test, trace);

      state.resultListMap[resultId] = resultList;
    },
  },
  extraReducers: builder => {
    builder.addCase(updateTestResult.fulfilled, (state, {payload: {resultList, resultId}}) => {
      state.resultListMap[resultId] = resultList;
    });
  },
});

export const {replace} = ResultListSlice.actions;
export default ResultListSlice.reducer;
