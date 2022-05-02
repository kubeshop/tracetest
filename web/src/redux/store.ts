import {configureStore, ThunkAction, Action} from '@reduxjs/toolkit';
import TestAPI from 'redux/apis/Test.api';
import ResultListSlice from 'redux/slices/ResultList.slice';

export const store = configureStore({
  reducer: {
    [TestAPI.reducerPath]: TestAPI.reducer,
    resultList: ResultListSlice,
  },
  middleware: getDefaultMiddleware => getDefaultMiddleware().concat(TestAPI.middleware),
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType, RootState, unknown, Action<string>>;
