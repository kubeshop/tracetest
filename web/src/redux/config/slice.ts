import {createSlice, PayloadAction} from '@reduxjs/toolkit';

interface IConfigState {
  isDataStoreConfigured: boolean;
}

const initialState: IConfigState = {
  isDataStoreConfigured: true,
};

const configSlice = createSlice({
  name: 'config',
  initialState,
  reducers: {
    configInitiated(state, action: PayloadAction<{}>) {},
  },
});

export const {configInitiated} = configSlice.actions;

export default configSlice.reducer;
