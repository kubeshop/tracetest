import faker from '@faker-js/faker';
import {render} from '@testing-library/react';
import SpanMock from '../../../models/__mocks__/Span.mock';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import SpanDetail from '../SpanDetail';

test('Layout', () => {
  const {getByText} = render(
    <ReduxWrapperProvider>
      <SpanDetail span={SpanMock.model()} />
    </ReduxWrapperProvider>
  );
  expect(getByText('Attribute list')).toBeTruthy();
});
