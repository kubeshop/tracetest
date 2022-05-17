import {render} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import FailedTrace from '../index';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';

test('FailedTrace', () => {
  const result = render(
    <MemoryRouter>
      <ReduxWrapperProvider>
        <FailedTrace onRunTest={jest.fn()} testId="234" isDisplayingError onEdit={jest.fn()} />
      </ReduxWrapperProvider>
    </MemoryRouter>
  );
  expect(result.container).toMatchSnapshot();
});
