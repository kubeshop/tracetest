import {render, waitFor} from '@testing-library/react';
import TestResults from '../TestResults';
import {TestingModels} from '../../../utils/TestingModels';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';

test('SpanAttributesTable', async () => {
  const onSpanSelected = jest.fn();
  const onPointerDown = jest.fn();
  const onHeaderClick = jest.fn();
  const {getByText} = render(
    <ReduxWrapperProvider>
      <TestResults
        result={TestingModels.testRunResult}
        visiblePortion={0}
        onSpanSelected={onSpanSelected}
        onPointerDown={onPointerDown}
        onHeaderClick={onHeaderClick}
      />
    </ReduxWrapperProvider>
  );
  await waitFor(() => getByText('Trace Information'));
});
