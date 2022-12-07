import {createSlice} from '@reduxjs/toolkit';
import {ISetupConfigState, TSetupConfigSliceActions} from 'types/Config.types';

export const initialState: ISetupConfigState = {
  draftConfig: {},
  isFormValid: true,
};

const setupConfigSlice = createSlice<ISetupConfigState, TSetupConfigSliceActions, 'setupConfig'>({
  name: 'setupConfig',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    setIsFormValid(state, {payload: {isValid}}) {
      state.isFormValid = isValid;
    },
    setDraftConfig(state, {payload: {draftConfig}}) {
      state.draftConfig = {
        ...state.draftConfig,
        ...draftConfig,
      };
    },
  },
});

export const {reset, setIsFormValid, setDraftConfig} = setupConfigSlice.actions;
export default setupConfigSlice.reducer;
