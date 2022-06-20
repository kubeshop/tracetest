import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {ISpanState, TSpan} from '../../types/Span.types';
import {setSelectedAssertion} from './TestDefinition.slice';

export const initialState: ISpanState = {
  affectedSpans: [],
  focusedSpan: '',
  selectedSpan: undefined,
};

const testDefinitionSlice = createSlice({
  name: 'spans',
  initialState,
  reducers: {
    setAffectedSpans(state, {payload: {spanIds}}: PayloadAction<{spanIds: string[]}>) {
      state.affectedSpans = spanIds;
      state.focusedSpan = spanIds[0] || '';
    },
    setSelectedSpan(state, {payload: {span}}: PayloadAction<{span: TSpan}>) {
      state.selectedSpan = span;
    },
    setFocusedSpan(state, {payload: {spanId}}: PayloadAction<{spanId: string}>) {
      state.focusedSpan = spanId;
    },
    clearAffectedSpans(state) {
      state.affectedSpans = [];
      state.focusedSpan = '';
    },
  },
  extraReducers: builder => {
    builder.addCase(setSelectedAssertion, (state, {payload: assertionResult}) => {
      state.affectedSpans = assertionResult?.spanIds ?? [];
      state.focusedSpan = state.affectedSpans[0] || '';
    });
  },
});

export const {clearAffectedSpans, setAffectedSpans, setSelectedSpan, setFocusedSpan} = testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
