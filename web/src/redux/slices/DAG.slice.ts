import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {Node, Edge, NodeChange, applyNodeChanges} from 'react-flow-renderer';

import DAGModel from 'models/DAG.model';
import {TSpan} from 'types/Span.types';
import {clearAffectedSpans, setAffectedSpans, setMatchedSpans, setSelectedSpan} from './Span.slice';
import {setSelectedAssertion} from './TestDefinition.slice';

interface IState {
  edges: Edge[];
  nodes: Node[];
}

const initialState: IState = {
  edges: [],
  nodes: [],
};

const dagSlice = createSlice({
  name: 'dag',
  initialState,
  reducers: {
    initNodes(state, {payload}: PayloadAction<{spans: TSpan[]}>) {
      const {edges, nodes} = DAGModel(payload.spans);
      state.edges = edges;
      state.nodes = nodes;
    },
    onNodesChange(state, {payload}: PayloadAction<{changes: NodeChange[]}>) {
      state.nodes = applyNodeChanges(payload.changes, state.nodes);
    },
  },
  extraReducers: builder => {
    builder
      .addCase(clearAffectedSpans, state => {
        state.nodes = state.nodes.map(node => {
          return {...node, data: {...node.data, isAffected: false}};
        });
      })
      .addCase(setAffectedSpans, (state, {payload: {spanIds}}) => {
        state.nodes = state.nodes.map(node => {
          const isAffected = spanIds.includes(node.id);
          return {...node, data: {...node.data, isAffected}};
        });
      })
      .addCase(setMatchedSpans, (state, {payload: {spanIds}}) => {
        state.nodes = state.nodes.map(node => {
          const isMatched = spanIds.includes(node.id);
          return {...node, data: {...node.data, isMatched}};
        });
      })
      .addCase(setSelectedAssertion, (state, {payload: assertionResult}) => {
        const spanIds = assertionResult?.spanIds ?? [];
        state.nodes = state.nodes.map(node => {
          const isAffected = spanIds.includes(node.id);
          return {...node, data: {...node.data, isAffected}};
        });
      })
      .addCase(setSelectedSpan, (state, {payload: {span}}) => {
        state.nodes = state.nodes.map(node => {
          const selected = span.id === node.id;
          return {...node, selected};
        });
      });
  },
});

export const {initNodes, onNodesChange} = dagSlice.actions;
export default dagSlice.reducer;
