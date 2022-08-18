import {createListenerMiddleware} from '@reduxjs/toolkit';
import {LOCATION_CHANGE} from 'redux-first-history';
import {parse} from 'query-string';
import RouterActions from './actions/Router.actions';
import {RootState} from './store';

const listener = createListenerMiddleware();

const sideEffectActionList = [RouterActions.updateSelectedAssertion, RouterActions.updateSelectedSpan];

const RouterMiddleware = () => ({
  middleware: listener.middleware,
  startListening(params = {}) {
    return listener.startListening({
      predicate: ({type = ''}) =>
        type === LOCATION_CHANGE || (type.endsWith('/fulfilled') && !type.startsWith('router/')),
      effect(_, {dispatch, getState}) {
        const {
          router: {location},
        } = getState() as RootState;

        const search = parse(location?.search || '');

        sideEffectActionList.forEach(action => {
          dispatch(action({search, params}));
        });
      },
    });
  },
});

export default RouterMiddleware();
