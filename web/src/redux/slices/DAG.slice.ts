import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {Node, Edge, NodeChange, applyNodeChanges} from 'react-flow-renderer';

import {endpoints} from 'redux/apis/TraceTest.api';
import DAGService from 'services/DAG.service';
import SpanService from 'services/Span.service';
import {clearAffectedSpans, setAffectedSpans, setSelectedSpan, setMatchedSpans} from './Span.slice';
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
    reset() {
      return initialState;
    },
    onNodesChange(state, {payload}: PayloadAction<{changes: NodeChange[]}>) {
      state.nodes = applyNodeChanges(payload.changes, state.nodes);
    },
  },
  extraReducers: builder => {
    builder
      .addCase(clearAffectedSpans, state => {
        console.log('### clearAffectedSpans');
        const newNodes = state.nodes.map(node => {
          return {...node, data: {...node.data, isAffected: false}};
        });
        state.nodes = newNodes;
      })
      .addCase(setAffectedSpans, (state, {payload: {spanIds}}) => {
        console.log('### setAffectedSpans', spanIds);
        const newNodes = state.nodes.map(node => {
          const isAffected = spanIds.includes(node.id);
          return {...node, data: {...node.data, isAffected}};
        });
        state.nodes = newNodes;
      })
      .addCase(setSelectedAssertion, (state, {payload: assertionResult}) => {
        console.log('### setSelectedAssertion', assertionResult);
        const spanIds = assertionResult?.spanIds ?? [];
        const newNodes = state.nodes.map(node => {
          const isAffected = spanIds.includes(node.id);
          return {...node, data: {...node.data, isAffected}};
        });
        state.nodes = newNodes;
      })
      .addCase(setSelectedSpan, (state, {payload: {span}}) => {
        console.log('### setSelectedSpan', span);
        const newNodes = state.nodes.map(node => {
          const selected = span.id === node.id;
          return {...node, selected};
        });
        state.nodes = newNodes;
      })
      .addCase(setMatchedSpans, (state, {payload: {spanIds}}) => {
        console.log('### setMatchedSpans', spanIds);
        const newNodes = state.nodes.map(node => {
          const isMatched = spanIds.includes(node.id);
          return {...node, data: {...node.data, isMatched}};
        });
        state.nodes = newNodes;
      })
      .addMatcher(endpoints.getRunById.matchFulfilled, (state, {payload}) => {
        console.log('### getRunById.matchFulfilled', payload);
        const spans = payload?.trace?.spans ?? [];
        const nodeList = SpanService.getNodeListFromSpanList(spans);
        const {edges, nodes} = DAGService.getNodesAndEdges(nodeList);
        state.edges = edges;
        state.nodes = nodes;
      })
      .addMatcher(
        action => action.type === 'tests/queries/queryResultPatched',
        (state, {payload}) => {
          console.log('### getRunById.tests/queries/queryResultPatched', payload);
          if (payload?.queryCacheKey?.startsWith('getRunById')) {
            const spans = payload?.patches?.[0]?.value?.trace?.spans ?? [];
            const nodeList = SpanService.getNodeListFromSpanList(spans);
            const {edges, nodes} = DAGService.getNodesAndEdges(nodeList);
            state.edges = edges;
            state.nodes = nodes;
          }
        }
      );
  },
});

export const {reset, onNodesChange} = dagSlice.actions;
export default dagSlice.reducer;
