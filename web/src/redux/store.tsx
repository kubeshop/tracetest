import {Action, configureStore, Middleware, ThunkAction} from '@reduxjs/toolkit';
import {createReduxHistoryContext} from 'redux-first-history';
import {createBrowserHistory} from 'history';
import TestAPI from 'redux/apis/TraceTest.api';
import TestSpecs from 'redux/slices/TestSpecs.slice';
import Spans from 'redux/slices/Span.slice';
import CreateTest from 'redux/slices/CreateTest.slice';
import DAG from 'redux/slices/DAG.slice';
import Trace from 'redux/slices/Trace.slice';
import CreateTransaction from 'redux/slices/CreateTransaction.slice';
import User from 'redux/slices/User.slice';
import RouterMiddleware from './Router.middleware';
import OtelRepoApi from './apis/OtelRepo.api';
import TestOutputs from './testOutputs/slice';

const {createReduxHistory, routerMiddleware, routerReducer} = createReduxHistoryContext({
  history: createBrowserHistory(),
});

export const middlewares: Middleware[] = [TestAPI.middleware, OtelRepoApi.middleware];

export const reducers = {
  [TestAPI.reducerPath]: TestAPI.reducer,
  [OtelRepoApi.reducerPath]: OtelRepoApi.reducer,

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
