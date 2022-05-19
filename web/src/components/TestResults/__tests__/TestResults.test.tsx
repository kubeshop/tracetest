import {render} from '@testing-library/react';
import TestResults from '../TestResults';
import {TestingModels} from '../../../utils/TestingModels';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';

test('TestResults', async () => {
  const onSpanSelected = jest.fn();
  const {getAllByTestId} = render(
    <ReduxWrapperProvider>
      <TestResults assertionResultList={[]} result={TestingModels.testRunResult} onSelectSpan={onSpanSelected} />
    </ReduxWrapperProvider>
  );

  expect(getAllByTestId('assertion-card-list').length).toBe(1);
});
