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
import CreateTransaction from 'redux/slices/CreateTransaction.slice';
import User from 'redux/slices/User.slice';
import TestOutputs from 'redux/testOutputs/slice';
import RouterMiddleware from './Router.middleware';

const {createReduxHistory, routerMiddleware, routerReducer} = createReduxHistoryContext({
  history: createBrowserHistory(),
});

export const middlewares: Middleware[] = [TracetestAPI.middleware, OtelRepoAPI.middleware];

export const reducers = {
  [TracetestAPI.reducerPath]: TracetestAPI.reducer,
  [OtelRepoAPI.reducerPath]: OtelRepoAPI.reducer,

  spans: Spans,
  dag: DAG,
  trace: Trace,
  testSpecs: TestSpecs,
  createTest: CreateTest,
  createTransaction: CreateTransaction,
  user: User,
  testOutputs: TestOutputs,
};

export const store = configureStore({
  reducer: {
    ...reducers,
    router: routerReducer,
  },
  middleware: getDefaultMiddleware =>
    getDefaultMiddleware()
      .prepend(RouterMiddleware.middleware)
      .concat(...middlewares, routerMiddleware),
});

export const history = createReduxHistory(store);

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
