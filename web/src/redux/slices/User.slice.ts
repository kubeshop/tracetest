import {PayloadAction, createAction, createSlice} from '@reduxjs/toolkit';
import UserPreferencesService from 'services/UserPreferences.service';
import {IUserState, TUserPreferenceKey, TUserPreferenceValue} from 'types/User.types';

export const initialState: IUserState = {
  preferences: UserPreferencesService.get(),
  runOriginPath: '/',
};

interface ISetUserPreferencesProps {
  key: TUserPreferenceKey;
  value: TUserPreferenceValue;
}

export const setUserPreference = createAction('user/setUserPreference', ({key, value}: ISetUserPreferencesProps) => {
  return {
    payload: {preferences: UserPreferencesService.set(key, value)},
  };
});

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    runOriginPathAdded(state, {payload}: PayloadAction<string>) {
      state.runOriginPath = payload;
    },
  },
  extraReducers: builder => {
    builder.addCase(setUserPreference, (state, {payload: {preferences}}) => {
      state.preferences = preferences;
    });
  },
});

export const {runOriginPathAdded} = userSlice.actions;
export default userSlice.reducer;
