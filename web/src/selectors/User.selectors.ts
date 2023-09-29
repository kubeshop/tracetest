import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {IUserPreferences, TUserPreferenceKey} from 'types/User.types';

const stateSelector = (state: RootState) => state;
const preferenceKeySelector = (state: RootState, key: TUserPreferenceKey) => key;

const UserSelectors = () => ({
  selectUserPreference: createSelector(
    stateSelector,
    preferenceKeySelector,
    <K extends keyof IUserPreferences>({user}: RootState, key: K): IUserPreferences[K] => {
      return user.preferences[key];
    }
  ),

  selectRunOriginPath: (state: RootState) => state.user.runOriginPath,
});

export default UserSelectors();
