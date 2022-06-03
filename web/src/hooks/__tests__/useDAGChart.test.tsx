import {renderHook} from '@testing-library/react-hooks';
import {skeletonNodeList} from '../../constants/Diagram.constants';
import {useDAGChart} from '../useDAGChart';

test('useDAGChart with empty span', () => {
  const {result} = renderHook(() => useDAGChart([]));
  expect(result.current).toEqual({});
});

test('useDAGChart with filled span', () => {
  const {result} = renderHook(() => useDAGChart(skeletonNodeList));
  expect(result.current).toEqual({
    dag: {
      data: {id: '1', parentIds: [], data: {}},
      dataChildren: [
        {
          child: {
            data: {id: '2', parentIds: ['1'], data: {}},
            dataChildren: [
              {
                child: {data: {id: '4', parentIds: ['2'], data: {}}, dataChildren: [], value: 2, x: 300, y: 375},
                points: [
                  {x: 500, y: 225},
                  {x: 300, y: 375},
                ],
              },
            ],
            value: 1,
            x: 500,
            y: 225,
          },
          points: [
            {x: 300, y: 75},
            {x: 500, y: 225},
          ],
        },
        {
          child: {data: {id: '3', parentIds: ['1'], data: {}}, dataChildren: [], value: 1, x: 300, y: 225},
          points: [
            {x: 300, y: 75},
            {x: 300, y: 225},
          ],
        },
        {
          child: {data: {id: '5', parentIds: ['1'], data: {}}, dataChildren: [], value: 1, x: 100, y: 225},
          points: [
            {x: 300, y: 75},
            {x: 100, y: 225},
          ],
        },
      ],
      value: 0,
      x: 300,
      y: 75,
    },
    layout: {width: 600, height: 450},
  });
});
