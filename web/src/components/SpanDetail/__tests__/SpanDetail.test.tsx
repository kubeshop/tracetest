import faker from '@faker-js/faker';
import {render} from '@testing-library/react';
import SpanMock from '../../../models/__mocks__/Span.mock';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import SpanDetail from '../SpanDetail';

test('Layout', () => {
  const result = render(
    <ReduxWrapperProvider>
      <SpanDetail resultId={faker.datatype.uuid()} testId={faker.datatype.uuid()} span={SpanMock.model()} />
    </ReduxWrapperProvider>
  );
  expect(result.container).toMatchSnapshot();
});
