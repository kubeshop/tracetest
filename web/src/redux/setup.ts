import {Middleware} from '@reduxjs/toolkit';

import OtelRepoAPI from 'redux/apis/OtelRepo';
import TestSpecs from 'redux/slices/TestSpecs.slice';
import Spans from 'redux/slices/Span.slice';
import CreateTest from 'redux/slices/CreateTest.slice';
import DAG from 'redux/slices/DAG.slice';
import Trace from 'redux/slices/Trace.slice';
import CreateTestSuite from 'redux/slices/CreateTestSuite.slice';
import User from 'redux/slices/User.slice';
import TestOutputs from 'redux/testOutputs/slice';

export const middlewares: Middleware[] = [OtelRepoAPI.middleware];

export const reducers = {
  [OtelRepoAPI.reducerPath]: OtelRepoAPI.reducer,

  spans: Spans,
  dag: DAG,
  trace: Trace,
  testSpecs: TestSpecs,
  createTest: CreateTest,
  createTestSuite: CreateTestSuite,
  user: User,
  testOutputs: TestOutputs,
};
