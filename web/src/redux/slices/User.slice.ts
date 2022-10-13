import {createAction, createSlice} from '@reduxjs/toolkit';
import UserPreferencesService from 'services/UserPreferences.service';
import {IUserState, TUserPreferenceKey, TUserPreferenceValue} from 'types/User.types';

export const initialState: IUserState = {
  preferences: UserPreferencesService.get(),
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

const testDefinitionSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {},
  extraReducers: builder => {
    builder.addCase(setUserPreference, (state, {payload: {preferences}}) => {
      state.preferences = preferences;
    });
  },
});

// export const {} = testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
