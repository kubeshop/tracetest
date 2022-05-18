import {render} from '@testing-library/react';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import {TestingModels} from '../../../utils/TestingModels';
import SpanDetail from '../SpanDetail';

test('Layout', () => {
  const result = render(
    <ReduxWrapperProvider>
      <SpanDetail resultId={TestingModels.resultId} testId={TestingModels.testId} span={TestingModels.span} />
    </ReduxWrapperProvider>
  );
  expect(result.container).toMatchSnapshot();
});
