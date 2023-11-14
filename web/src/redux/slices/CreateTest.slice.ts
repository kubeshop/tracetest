import {createSlice} from '@reduxjs/toolkit';
import {ICreateTestState, TCreateTestSliceActions} from 'types/Test.types';
import {SupportedPlugins} from 'constants/Common.constants';

export const initialState: ICreateTestState = {
  draftTest: {
    name: 'Untitled Test',
  },
  pluginName: SupportedPlugins.REST,
  isFormValid: true,
};

const createTestSlice = createSlice<ICreateTestState, TCreateTestSliceActions, 'createTest'>({
  name: 'createTest',
  initialState,
  reducers: {
    reset() {
      return initialState;
    },
    setIsFormValid(state, {payload: {isValid}}) {
      state.isFormValid = isValid;
    },
    setPlugin(
      state,
      {
        payload: {
          plugin: {name},
        },
      }
    ) {
      state.pluginName = name;
      state.draftTest = {};
    },
    setDraftTest(state, {payload: {draftTest}}) {
      state.draftTest = {
        ...state.draftTest,
        ...draftTest,
      };
    },
  },
});

export const {setPlugin, setDraftTest, reset, setIsFormValid} = createTestSlice.actions;
export default createTestSlice.reducer;
