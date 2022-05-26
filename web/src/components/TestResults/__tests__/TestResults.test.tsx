import {render} from '@testing-library/react';
import TestResults from '../TestResults';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';

test('TestResults', async () => {
  const onSpanSelected = jest.fn();
  const {getAllByTestId} = render(
    <ReduxWrapperProvider>
      <TestResults assertionResults={TestRunMock.model().result} onSelectSpan={onSpanSelected} testId="12345" />
    </ReduxWrapperProvider>
  );

  expect(getAllByTestId('assertion-card-list').length).toBe(1);
});
