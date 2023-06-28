import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {applyNodeChanges, Edge, MarkerType, Node, NodeChange} from 'react-flow-renderer';

import {NodeTypesEnum} from 'constants/DAG.constants';
import {theme} from 'constants/Theme.constants';
import DAGModel from 'models/DAG.model';
import Span from 'models/Span.model';
import {clearMatchedSpans, setMatchedSpans, setSelectedSpan} from './Span.slice';
import {setSelectedSpec} from './TestSpecs.slice';
import {outputsSelectedOutputsChanged} from '../testOutputs/slice';

export interface IDagState {
  edges: Edge[];
  nodes: Node[];
}

const initialState: IDagState = {
  edges: [],
  nodes: [],
};

const dagSlice = createSlice({
  name: 'dag',
  initialState,
  reducers: {
    initNodes(state, {payload}: PayloadAction<{spans: Span[]}>) {
      const {edges, nodes} = DAGModel(payload.spans, NodeTypesEnum.TestSpan);
      state.edges = edges;
      state.nodes = nodes;
    },
    onNodesChange(state, {payload}: PayloadAction<{changes: NodeChange[]}>) {
      state.nodes = applyNodeChanges(payload.changes, state.nodes);
    },
  },
  extraReducers: builder => {
    builder
      .addCase(clearMatchedSpans, state => {
        state.nodes = state.nodes.map(node => {
          return {...node, data: {...node.data, isMatched: false}};
        });
      })
      .addCase(setMatchedSpans, (state, {payload: {spanIds}}) => {
        state.nodes = state.nodes.map(node => {
          const isMatched = spanIds.includes(node.id);
          return {...node, data: {...node.data, isMatched}};
        });
      })
      .addCase(setSelectedSpec, (state, {payload: assertionResult}) => {
        const spanIds = assertionResult?.spanIds ?? [];
        state.nodes = state.nodes.map(node => {
          const isMatched = spanIds.includes(node.id);
          return {...node, data: {...node.data, isMatched}};
        });
      })
      .addCase(setSelectedSpan, (state, {payload: {span}}) => {
        state.edges = state.edges.map(edge => {
          const selected = span.id === edge.source;
          return {
            ...edge,
            animated: selected,
            markerEnd: {color: selected ? theme.color.interactive : theme.color.border, type: MarkerType.ArrowClosed},
            style: {stroke: selected ? theme.color.interactive : theme.color.border},
          };
        });

        state.nodes = state.nodes.map(node => {
          const selected = span.id === node.id;
          return {...node, selected};
        });
      })
      .addCase(outputsSelectedOutputsChanged, (state, {payload: outputs}) => {
        const spanIds = outputs.map(output => output.spanId);
        state.nodes = state.nodes.map(node => {
          const isMatched = spanIds.includes(node.id);
          return {...node, data: {...node.data, isMatched}};
        });
      });
  },
});

export const {initNodes, onNodesChange} = dagSlice.actions;
export default dagSlice.reducer;
