import {createSelector} from '@reduxjs/toolkit';
import {RootState} from 'redux/store';

const dagStateSelector = (state: RootState) => state.dag;

const DAGSelectors = () => ({
  selectEdges: createSelector(dagStateSelector, ({edges}) => edges),
  selectNodes: createSelector(dagStateSelector, ({nodes}) => nodes),
});

export default DAGSelectors();
