import {render} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import SpanMock from '../../../models/__mocks__/Span.mock';
import TraceNode from '../TraceNode';

test('TraceNode', () => {
  const {getByText} = render(
    <ReactFlowProvider>
      <TraceNode
        id="sdfkjn"
        type="ksdjfn"
        data={SpanMock.model()}
        dragHandle="sdkjf"
        isDragging={false}
        selected={false}
        isConnectable={false}
      />
    </ReactFlowProvider>
  );
  expect(getByText('General')).toBeTruthy();
});
