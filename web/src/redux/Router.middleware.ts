import {createListenerMiddleware} from '@reduxjs/toolkit';
import type {TypedStartListening} from '@reduxjs/toolkit';
import {LOCATION_CHANGE} from 'redux-first-history';
import {parse} from 'query-string';
import RouterActions from './actions/Router.actions';
import {runOriginPathAdded} from './slices/User.slice';
import {AppDispatch, RootState} from './store';

type AppStartListening = TypedStartListening<RootState, AppDispatch>;
const listener = createListenerMiddleware();
const startAppListening = listener.startListening as AppStartListening;

const sideEffectActionList = [RouterActions.updateSelectedAssertion, RouterActions.updateSelectedSpan];

const runUrlRegex = /^\/(test|testsuite)\/([^\/]+)\/run\/([^\/]+)(.*)$/;

const RouterMiddleware = () => ({
  middleware: listener.middleware,

  startListening(params = {}) {
    return startAppListening({
      predicate: ({type = ''}) =>
        type === LOCATION_CHANGE || (type.endsWith('/fulfilled') && !type.startsWith('router/')),
      effect(_, {dispatch, getState}) {
        const {
          router: {location},
        } = getState();

        const search = parse(location?.search || '');

        sideEffectActionList.forEach(action => {
          dispatch(action({search, params}));
        });
      },
    });
  },

  startListeningForLocationChange() {
    return startAppListening({
      predicate: action => {
        const pathname = action?.payload?.location?.pathname ?? '';
        return action?.type === LOCATION_CHANGE && pathname.match(runUrlRegex);
      },
      effect: async (_, {dispatch, getOriginalState, getState}) => {
        const {
          router: {location: prevLocation},
        } = getOriginalState();

        const {
          router: {location: currLocation},
        } = getState();

        const prevPathname = prevLocation?.pathname ?? '';

        if (!prevPathname.match(runUrlRegex) && !prevPathname.includes('/create')) {
          const defaultPath = currLocation?.pathname?.includes('testsuite') ? '/testsuites' : '/';
          dispatch(runOriginPathAdded(prevLocation?.pathname ?? defaultPath));
        }
      },
    });
  },
});

export default RouterMiddleware();
