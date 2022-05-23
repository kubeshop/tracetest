import {renderHook} from '@testing-library/react-hooks';
import {SemanticGroupNames} from '../../constants/SemanticGroupNames.constants';
import {useDAGChart} from '../useDAGChart';

test('useDAGChart with empty span', () => {
  const {result} = renderHook(() => useDAGChart({}));
  expect(result.current).toBe(undefined);
});

test('useDAGChart with filled span', () => {
  const {result} = renderHook(() =>
    useDAGChart({
      cesco: {
        id: '1',
        data: {
          attributes: {},
          type: SemanticGroupNames.Http,
          duration: 10,
          signature: [],
          attributeList: [],
          startTime: '',
          name: '',
          parentId: '',
          id: '',
          endTime: '',
        },
        parentIds: [],
      },
    })
  );
  expect(result.current).toEqual({
    dag: {data: {id: '1', parentIds: []}, dataChildren: [], value: 0, x: 100, y: 75},
    layout: {height: 150, width: 200},
  });
});
