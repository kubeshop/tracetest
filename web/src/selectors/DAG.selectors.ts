import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';
import {IDagState} from 'redux/slices/DAG.slice';

const dagStateSelector = (state: RootState): IDagState => state.dag;

const DAGSelectors = () => ({
  selectEdges: createSelector(dagStateSelector, ({edges}) => edges),
  selectNodes: createSelector(dagStateSelector, ({nodes}) => nodes),
});

export default DAGSelectors();
