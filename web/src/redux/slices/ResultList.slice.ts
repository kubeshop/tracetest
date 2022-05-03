import {createAsyncThunk, createSlice, PayloadAction} from '@reduxjs/toolkit';
import TestRunResultGateway from '../../gateways/TestRunResult.gateway';
import TraceService from '../../services/Trace.service';
import {TAssertionResult, TAssertionResultList} from '../../types/Assertion.types';
import {TTest} from '../../types/Test.types';
import {TTrace} from '../../types/Trace.types';

interface TTestResultListState {
  resultListMap: Record<string, TAssertionResult[]>;
}

interface TTestResultListReplaceParams {
  resultList: TAssertionResult[];
  resultId: string;
}

const initialState: TTestResultListState = {
  resultListMap: {},
};

export const updateTestResult = createAsyncThunk<
  TTestResultListReplaceParams,
  {trace: TTrace; resultId: string; test: TTest}
>('resultList/load', async ({trace, test, resultId}, {dispatch}) => {
  const resultList = TraceService.runTest(trace, test);

  await dispatch(
    TestRunResultGateway.update(test.testId || '', resultId, TraceService.parseAssertionResultListToTestResult(resultList))
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
      }: PayloadAction<{assertionResult: TAssertionResultList; test: TTest; trace: TTrace; resultId: string}>
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
