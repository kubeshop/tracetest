import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {TUserPreferenceKey} from 'types/User.types';

const stateSelector = (state: RootState) => state;
const preferenceKeySelector = (state: RootState, key: TUserPreferenceKey) => key;

const UserSelectors = () => ({
  selectUserPreference: createSelector(stateSelector, preferenceKeySelector, ({user}, key) => {
    return user.preferences[key];
  }),
});

export default UserSelectors();
