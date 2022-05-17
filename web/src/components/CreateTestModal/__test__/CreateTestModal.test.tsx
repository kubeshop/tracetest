import {render} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import CreateTestModal from '../CreateTestModal';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';

test('CreateTestModal', () => {
  const result = render(
    <MemoryRouter>
      <ReduxWrapperProvider>
        <CreateTestModal visible onClose={jest.fn()} />
      </ReduxWrapperProvider>
    </MemoryRouter>
  );
  expect(result.container).toMatchSnapshot();
});
