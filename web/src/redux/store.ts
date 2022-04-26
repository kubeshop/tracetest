import {configureStore, ThunkAction, Action} from '@reduxjs/toolkit';
import {testAPI} from 'redux/services/TestService';

export const store = configureStore({
  reducer: {
    [testAPI.reducerPath]: testAPI.reducer,
  },
  middleware: getDefaultMiddleware => getDefaultMiddleware().concat(testAPI.middleware),
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
