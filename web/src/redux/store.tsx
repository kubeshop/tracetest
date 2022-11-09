import {combineReducers} from 'redux';
import {Action, configureStore, ThunkAction} from '@reduxjs/toolkit';
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
import createWebSocketMiddleware from './WebSocket.middleware';

const {createReduxHistory, routerMiddleware, routerReducer} = createReduxHistoryContext({
  history: createBrowserHistory(),
});

const rootReducer = combineReducers({
  [TestAPI.reducerPath]: TestAPI.reducer,
  [OtelRepoApi.reducerPath]: OtelRepoApi.reducer,

  router: routerReducer,
  spans: Spans,
  dag: DAG,
  trace: Trace,
  testSpecs: TestSpecs,
  createTest: CreateTest,
  createTransaction: CreateTransaction,
  user: User,
  testOutputs: TestOutputs,
});

export const store = configureStore({
  reducer: rootReducer,
  middleware: getDefaultMiddleware =>
    getDefaultMiddleware()
      .prepend(RouterMiddleware.middleware)
      .concat(TestAPI.middleware, routerMiddleware, OtelRepoApi.middleware, createWebSocketMiddleware()),
});

export const history = createReduxHistory(store);

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof rootReducer>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
