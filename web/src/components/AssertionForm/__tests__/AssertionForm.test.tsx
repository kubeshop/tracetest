import {render} from '@testing-library/react';
import AssertionForm from '../AssertionForm';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';

const defaultProps = {
  onSubmit: jest.fn(),
  onCancel: jest.fn(),
  testId: 'testId',
  runId: 'runId',
};

describe('AssertionForm', () => {
  it('should render correctly', () => {
    fetchMock.mockResponse(JSON.stringify(['spanId']));
    const {container} = render(<AssertionForm {...defaultProps} />, {wrapper: ReduxWrapperProvider}); 

    expect(container.querySelector('form')).toBeInTheDocument();
  });
});
