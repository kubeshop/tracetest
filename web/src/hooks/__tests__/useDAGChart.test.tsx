import {renderHook} from '@testing-library/react-hooks';
import {skeletonNodeList} from '../../constants/Diagram.constants';
import {useDAGChart} from '../useDAGChart';

test('useDAGChart with empty span', () => {
  const {result} = renderHook(() => useDAGChart([]));
  expect(result.current).toEqual([]);
});

test('useDAGChart with filled span', () => {
  const {result} = renderHook(() => useDAGChart(skeletonNodeList));
  expect(result.current).toEqual([
    {id: '1', type: 'SkeletonNode', data: {}, position: {x: 300, y: 75}, sourcePosition: 'top'},
    {id: '5', type: 'SkeletonNode', data: {}, position: {x: 100, y: 225}, sourcePosition: 'top'},
    {id: '3', type: 'SkeletonNode', data: {}, position: {x: 300, y: 225}, sourcePosition: 'top'},
    {id: '2', type: 'SkeletonNode', data: {}, position: {x: 500, y: 225}, sourcePosition: 'top'},
    {id: '4', type: 'SkeletonNode', data: {}, position: {x: 300, y: 375}, sourcePosition: 'top'},
    {
      id: '1_2',
      source: '1',
      target: '2',
      data: {},
      labelShowBg: false,
      animated: false,
      arrowHeadType: 'arrowclosed',
      style: {stroke: '#C9CEDB'},
    },
    {
      id: '1_3',
      source: '1',
      target: '3',
      data: {},
      labelShowBg: false,
      animated: false,
      arrowHeadType: 'arrowclosed',
      style: {stroke: '#C9CEDB'},
    },
    {
      id: '1_5',
      source: '1',
      target: '5',
      data: {},
      labelShowBg: false,
      animated: false,
      arrowHeadType: 'arrowclosed',
      style: {stroke: '#C9CEDB'},
    },
    {
      id: '2_4',
      source: '2',
      target: '4',
      data: {},
      labelShowBg: false,
      animated: false,
      arrowHeadType: 'arrowclosed',
      style: {stroke: '#C9CEDB'},
    },
  ]);
});
