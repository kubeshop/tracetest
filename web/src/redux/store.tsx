import {Action, configureStore, ThunkAction} from '@reduxjs/toolkit';
import TestAPI from './apis/TraceTest.api';
import DAG from './slices/DAG.slice';
import SpanS from './slices/Span.slice';
import TestDefinition from './slices/TestDefinition.slice';

export const store = configureStore({
  reducer: {
    [TestAPI.reducerPath]: TestAPI.reducer,
    dag: DAG,
    spans: SpanS,
    testDefinition: TestDefinition,
  },
  middleware: getDefaultMiddleware => getDefaultMiddleware().concat(TestAPI.middleware),
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
