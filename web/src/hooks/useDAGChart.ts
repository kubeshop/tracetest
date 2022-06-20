import {Dag} from 'd3-dag';
import {noop} from 'lodash';
import {useCallback, useEffect, useState} from 'react';
import {TSpan} from 'types/Span.types';
import DiagramService, {IDAGNode, TElementList, TNode} from '../services/DAG.service';

export const useDAGChart = <T>(nodeList: TNode<T>[], selectedSpan?: TSpan, onSelectSpan = noop) => {
  const [dag, setDag] = useState<Dag<IDAGNode<unknown>, undefined>>();
  const [elementList, setElementList] = useState<TElementList>([]);

  const loadData = useCallback(() => {
    const info = DiagramService.getDagInfo(nodeList);

    setDag(info.dag);
    setElementList(info.elementList);
  }, [nodeList]);

  useEffect(() => {
    if (dag) {
      const [dagNode] = dag.descendants();
      const node = elementList.find(({id}) => id === dagNode?.data.id);

      if (!selectedSpan && node) {
        onSelectSpan(node.id);
      }
    }
  }, [dag, elementList, onSelectSpan, selectedSpan]);

  useEffect(() => {
    if (nodeList.length) loadData();
  }, [loadData, nodeList]);

  return elementList;
};
