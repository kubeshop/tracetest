import {Action, configureStore, Middleware, ThunkAction} from '@reduxjs/toolkit';
import {createReduxHistoryContext} from 'redux-first-history';
import {createBrowserHistory} from 'history';
import TracetestAPI from 'redux/apis/Tracetest';
import OtelRepoAPI from 'redux/apis/OtelRepo';
import TestSpecs from 'redux/slices/TestSpecs.slice';
import Spans from 'redux/slices/Span.slice';
import CreateTest from 'redux/slices/CreateTest.slice';
import DAG from 'redux/slices/DAG.slice';
import Trace from 'redux/slices/Trace.slice';
import CreateTestSuite from 'redux/slices/CreateTestSuite.slice';
import User from 'redux/slices/User.slice';
import TestOutputs from 'redux/testOutputs/slice';
import RouterMiddleware from './Router.middleware';

const {createReduxHistory, routerMiddleware, routerReducer} = createReduxHistoryContext({
  history: createBrowserHistory(),
});

TracetestAPI.create();

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

export const store = configureStore({
  reducer: {
    ...reducers,
    [TracetestAPI.instance.reducerPath]: TracetestAPI.instance.reducer,
    router: routerReducer,
  },
  middleware: getDefaultMiddleware =>
    getDefaultMiddleware()
      .prepend(RouterMiddleware.middleware)
      .concat(TracetestAPI.instance.middleware, ...middlewares, routerMiddleware),
});

export const history = createReduxHistory(store);

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
