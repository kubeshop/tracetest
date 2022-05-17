import {renderHook} from '@testing-library/react-hooks';
import {useOnDeleteCallback} from '../useOnDeleteCallback';

test('useOnDeleteCallback', async () => {
  const onDelete = jest.fn();
  const {result} = renderHook(() => useOnDeleteCallback(onDelete));

  result.current({domEvent: new MouseEvent('click') as any});
  expect(onDelete).toBeCalled();
});
