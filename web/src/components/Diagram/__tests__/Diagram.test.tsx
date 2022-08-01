import {render, waitFor} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {TestState} from '../../../constants/TestRun.constants';
import TraceMock from '../../../models/__mocks__/Trace.mock';
import Diagram, {SupportedDiagrams} from '../Diagram';

jest.mock('../../../services/Analytics/Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

describe.skip('Diagram', () => {
  test('Diagram DAG', async () => {
    const {getByText} = render(
      <div style={{width: 200, height: 200}}>
        <ReactFlowProvider>
          <Diagram runState={TestState.FINISHED} trace={TraceMock.model()} type={SupportedDiagrams.DAG} />
        </ReactFlowProvider>
      </div>
    );
    await waitFor(() => getByText('HTTP'));
  });

  test('Diagram Timeline', async () => {
    const {getByText} = render(
      <div style={{width: 200, height: 200}}>
        <ReactFlowProvider>
          <Diagram runState={TestState.FINISHED} trace={TraceMock.model()} type={SupportedDiagrams.Timeline} />
        </ReactFlowProvider>
      </div>
    );

    await waitFor(() => getByText('Duration (ms)'));
  });

  test('should render the skeleton', async () => {
    const {getByTestId} = render(
      <div style={{width: 200, height: 200}}>
        <ReactFlowProvider>
          <Diagram runState={TestState.AWAITING_TRACE} trace={TraceMock.model()} type={SupportedDiagrams.Timeline} />
        </ReactFlowProvider>
      </div>
    );

    expect(getByTestId('skeleton-diagram')).toBeInTheDocument();
  });
});
