import {render, waitFor} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import SpanMock from '../../../models/__mocks__/Span.mock';
import TraceMock from '../../../models/__mocks__/Trace.mock';
import Diagram, {SupportedDiagrams} from '../Diagram';

jest.mock('../../../services/Analytics/Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

test('Diagram DAG', async () => {
  const {getByText} = render(
    <div style={{width: 200, height: 200}}>
      <ReactFlowProvider>
        <Diagram
          affectedSpans={[]}
          type={SupportedDiagrams.DAG}
          trace={TraceMock.model()}
          onSelectSpan={jest.fn()}
          selectedSpan={SpanMock.model()}
        />
      </ReactFlowProvider>
    </div>
  );
  await waitFor(() => getByText('HTTP'));
});

test('Diagram Timeline', async () => {
  const {getByText} = render(
    <div style={{width: 200, height: 200}}>
      <ReactFlowProvider>
        <Diagram
          affectedSpans={[]}
          type={SupportedDiagrams.Timeline}
          trace={TraceMock.model()}
          onSelectSpan={jest.fn()}
          selectedSpan={SpanMock.model()}
        />
      </ReactFlowProvider>
    </div>
  );

  await waitFor(() => getByText('Duration (ms)'));
});
