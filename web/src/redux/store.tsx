import {Action, configureStore, ThunkAction} from '@reduxjs/toolkit';
import TestAPI from 'redux/apis/TraceTest.api';
import TestDefinition from 'redux/slices/TestDefinition.slice';
import Spans from 'redux/slices/Span.slice';
import CreateTest from 'redux/slices/CreateTest.slice';
import DAG from 'redux/slices/DAG.slice';

export const store = configureStore({
  reducer: {
    [TestAPI.reducerPath]: TestAPI.reducer,
    spans: Spans,
    dag: DAG,
    testDefinition: TestDefinition,
    createTest: CreateTest,
  },
  middleware: getDefaultMiddleware => getDefaultMiddleware().concat(TestAPI.middleware),
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
