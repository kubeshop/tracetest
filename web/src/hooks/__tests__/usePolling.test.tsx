import {renderHook} from '@testing-library/react-hooks';
import usePolling from '../usePolling';

test('usePolling', async () => {
  const callback = jest.fn();
  const payload = {callback, isPolling: true, delay: 0};
  await renderHook(() => usePolling(payload));

  expect(callback).toBeCalledTimes(0);
});
