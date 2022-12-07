import {createSelector} from '@reduxjs/toolkit';

import UserSelectors from 'selectors/User.selectors';
import {RootState} from '../store';

const stateSelector = (state: RootState) => state;

const configStateSelector = (state: RootState) => state.config;

export const selectIsDataStoreConfigured = createSelector(
  configStateSelector,
  ({isDataStoreConfigured}) => isDataStoreConfigured
);

export const selectInitConfigSetup = createSelector(stateSelector, state =>
  Boolean(UserSelectors.selectUserPreference(state, 'initConfigSetup'))
);

export const selectShouldDisplayConfigSetup = createSelector(
  selectIsDataStoreConfigured,
  selectInitConfigSetup,
  (isDataStoreConfigured, initConfigSetup) => !isDataStoreConfigured && initConfigSetup
);
