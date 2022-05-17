import {render} from '@testing-library/react';
import Http from '../components/Http/Http';
import {TestingModels} from '../../../utils/TestingModels';

test('Http', () => {
  const {getAllByTestId} = render(
    <Http
      testId={TestingModels.testId}
      resultId={TestingModels.resultId}
      onCreateAssertion={jest.fn()}
      assertionsResultList={[
        {
          assertion: TestingModels.assertion,
          assertionResultList: [TestingModels.spanAssertionResult],
        },
      ]}
      span={TestingModels.span}
    />
  );
  expect(getAllByTestId('assertion-check-property').length).toBe(1);
});
