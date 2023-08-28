import {Action, configureStore, ThunkAction} from '@reduxjs/toolkit';
import {createReduxHistoryContext} from 'redux-first-history';
import {createBrowserHistory} from 'history';
import TracetestAPI from 'redux/apis/Tracetest';

import RouterMiddleware from './Router.middleware';
import {middlewares, reducers} from './setup';

const {createReduxHistory, routerMiddleware, routerReducer} = createReduxHistoryContext({
  history: createBrowserHistory(),
});

TracetestAPI.create();

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
