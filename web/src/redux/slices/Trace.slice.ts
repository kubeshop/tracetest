import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {applyNodeChanges, Edge, MarkerType, Node, NodeChange} from 'react-flow-renderer';

import {NodeTypesEnum} from 'constants/DAG.constants';
import {theme} from 'constants/Theme.constants';
import DAGModel from 'models/DAG.model';
import Span from 'models/Span.model';

export interface ITraceState {
  edges: Edge[];
  matchedSpans: string[];
  nodes: Node[];
  searchText: string;
  selectedSpan: string;
}

const initialState: ITraceState = {
  edges: [],
  matchedSpans: [],
  nodes: [],
  searchText: '',
  selectedSpan: '',
};

const traceSlice = createSlice({
  name: 'trace',
  initialState,
  reducers: {
    initNodes(state, {payload}: PayloadAction<{spans: Span[]}>) {
      const {edges, nodes} = DAGModel(payload.spans, NodeTypesEnum.TraceSpan);
      state.edges = edges;
      state.nodes = nodes;
      // Clear state
      state.matchedSpans = [];
      state.searchText = '';
      state.selectedSpan = '';
    },
    changeNodes(state, {payload}: PayloadAction<{changes: NodeChange[]}>) {
      state.nodes = applyNodeChanges(payload.changes, state.nodes);
    },
    selectSpan(state, {payload}: PayloadAction<{spanId: string}>) {
      state.selectedSpan = payload.spanId;

      state.edges = state.edges.map(edge => {
        const selected = payload.spanId === edge.source;
        return {
          ...edge,
          animated: selected,
          markerEnd: {color: selected ? theme.color.interactive : theme.color.border, type: MarkerType.ArrowClosed},
          style: {stroke: selected ? theme.color.interactive : theme.color.border},
        };
      });

      state.nodes = state.nodes.map(node => {
        const selected = payload.spanId === node.id;
        return {...node, selected};
      });
    },
    matchSpans(state, {payload}: PayloadAction<{spanIds: string[]}>) {
      state.matchedSpans = payload.spanIds;

      state.nodes = state.nodes.map(node => {
        const isMatched = payload.spanIds.includes(node.id);
        return {...node, data: {...node.data, isMatched}};
      });
    },
    setSearchText(state, {payload}: PayloadAction<{searchText: string}>) {
      state.searchText = payload.searchText.toLowerCase();
    },
  },
});

export const {initNodes, changeNodes, selectSpan, matchSpans, setSearchText} = traceSlice.actions;
export default traceSlice.reducer;
