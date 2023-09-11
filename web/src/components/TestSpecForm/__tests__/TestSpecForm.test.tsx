import {render} from 'test-utils';
import {ReactFlowProvider} from 'react-flow-renderer';
import TestSpecForm from '../TestSpecForm';

const defaultProps = {
  onSubmit: jest.fn(),
  onCancel: jest.fn(),
  onClearSelectorSuggestions: jest.fn(),
  onClickPrevSelector: jest.fn(),
  testId: 'testId',
  runId: 1,
  isValid: false,
  onIsValid: jest.fn(),
};

describe('TestSpecForm', () => {
  it('should render correctly', () => {
    fetchMock.mockResponse(JSON.stringify(['spanId']));
    const {container} = render(
      <ReactFlowProvider>
        <TestSpecForm {...defaultProps} />
      </ReactFlowProvider>
    );

    expect(container.querySelector('form')).toBeInTheDocument();
  });
});
