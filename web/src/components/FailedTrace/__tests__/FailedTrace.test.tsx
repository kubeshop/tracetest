import {render} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import FailedTrace from '../index';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';

test('FailedTrace', () => {
  const test = TestMock.model();
  const {getByText} = render(
    <MemoryRouter>
      <ReduxWrapperProvider>
        <FailedTrace
          onRunTest={jest.fn()}
          testId={test.id}
          isDisplayingError
          onEdit={jest.fn()}
          run={TestRunMock.model()}
        />
      </ReduxWrapperProvider>
    </MemoryRouter>
  );
  expect(getByText('Rerun Test')).toBeTruthy();
});
