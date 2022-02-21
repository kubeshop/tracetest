import {configureStore, ThunkAction, Action} from '@reduxjs/toolkit';
import {testAPI} from '../services/TestService';

export const store = configureStore({
  reducer: {
    [testAPI.reducerPath]: testAPI.reducer,
  },
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
