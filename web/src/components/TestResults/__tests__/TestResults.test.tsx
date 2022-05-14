import {render} from '@testing-library/react';
import TestResults from '../TestResults';
import {TestingModels} from '../../../utils/TestingModels';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';

test('SpanAttributesTable', () => {
  const setMax = jest.fn();
  const onSpanSelected = jest.fn();
  const onPointerDown = jest.fn();
  const setHeight = jest.fn();
  const result = render(
    <ReduxWrapperProvider>
      <TestResults
        result={TestingModels.testRunResult}
        visiblePortion={0}
        max={0}
        setMax={setMax}
        onSpanSelected={onSpanSelected}
        height={100}
        onPointerDown={onPointerDown}
        setHeight={setHeight}
      />
    </ReduxWrapperProvider>
  );
  expect(result.container).toMatchSnapshot();
});
