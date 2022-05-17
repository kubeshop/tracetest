import {render} from '@testing-library/react';
import {Provider} from 'react-redux';
import {store} from '../../../redux/store';
import {TestingModels} from '../../../utils/TestingModels';
import AssertionsTable from '../index';

test('AssetionsTable', () => {
  const result = render(
    <Provider store={store}>
      <AssertionsTable
        assertionResults={[TestingModels.spanAssertionResult]}
        sort={1}
        assertion={TestingModels.assertion}
        span={TestingModels.span}
        resultId={TestingModels.resultId}
        testId={TestingModels.testId}
      />
    </Provider>
  );
  expect(result.container).toMatchSnapshot();
});
