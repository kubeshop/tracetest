import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {ISpanState} from 'types/Span.types';
import Span from 'models/Span.model';
import {outputsSelectedOutputsChanged} from '../testOutputs/slice';
import {setSelectedSpec} from './TestSpecs.slice';

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
    setSelectedSpan(state, {payload: {span}}: PayloadAction<{span: Span}>) {
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
    builder
      .addCase(setSelectedSpec, (state, {payload: assertionResult}) => {
        state.matchedSpans = assertionResult?.spanIds ?? [];
        state.focusedSpan = state.matchedSpans[0] || '';
      })
      .addCase(outputsSelectedOutputsChanged, (state, {payload: outputs}) => {
        state.matchedSpans = outputs.map(output => output.spanId);
        state.focusedSpan = state.matchedSpans[0] || '';
      });
  },
});

export const {clearMatchedSpans, setMatchedSpans, clearSelectedSpan, setSelectedSpan, setFocusedSpan} =
  testDefinitionSlice.actions;
export default testDefinitionSlice.reducer;
