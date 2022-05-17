import {renderHook} from '@testing-library/react-hooks';
import {TestingModels} from '../../../utils/TestingModels';
import {useParserAssertionListMemo} from '../useParserAssertionListMemo';

test('useParserAssertionListMemo', () => {
  const spanListAssertionResult = TestingModels.assertionResult;
  let spanListAssertionResult1 = spanListAssertionResult.spanListAssertionResult;
  const {result} = renderHook(() => useParserAssertionListMemo(spanListAssertionResult1, []));
  expect(result.current).toStrictEqual([
    {
      actualValue: 'New request',
      comparison: 'EQUALS',
      hasPassed: false,
      key: '-',
      property: '',
      spanId: '',
      spanLabels: ['', '#'],
      value: '',
    },
  ]);
});
