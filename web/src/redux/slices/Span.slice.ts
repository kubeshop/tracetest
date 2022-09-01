import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {ISpanState, TSpan} from 'types/Span.types';
import {setSelectedAssertion} from './TestDefinition.slice';

export const initialState: ISpanState = {
  focusedSpan: '',
  matchedSpans: [],
  selectedSpan: undefined,
};

const testDefinitionSlice = createSlice({
  name: 'spans',
  initialState,
  reducers: {
    setMatchedSpans(state, {payload: {spanIds}}: PayloadAction<{spanIds: string[]}>) {
      state.matchedSpans = spanIds;
      state.focusedSpan = spanIds[0] || '';
    },
    setSelectedSpan(state, {payload: {span}}: PayloadAction<{span: TSpan}>) {
      state.selectedSpan = span;
    },
    setFocusedSpan(state, {payload: {spanId}}: PayloadAction<{spanId: string}>) {
      state.focusedSpan = spanId;
    },
    clearMatchedSpans(state) {
      state.matchedSpans = [];
      state.focusedSpan = '';
    },
    clearSelectedSpan(state) {
      state.selectedSpan = undefined;
    },
  },
  extraReducers: builder => {
    builder.addCase(setSelectedAssertion, (state, {payload: assertionResult}) => {
      state.matchedSpans = assertionResult?.spanIds ?? [];
      state.focusedSpan = state.matchedSpans[0] || '';
    });
  },
});

export const {clearMatchedSpans, setMatchedSpans, clearSelectedSpan, setSelectedSpan, setFocusedSpan} =
  testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
