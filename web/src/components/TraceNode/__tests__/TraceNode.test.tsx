import {render} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {TestingModels} from '../../../utils/TestingModels';
import TraceNode from '../TraceNode';

test('TraceNode', () => {
  const result = render(
    <ReactFlowProvider>
      <TraceNode
        id="sdfkjn"
        type="ksdjfn"
        data={TestingModels.span}
        dragHandle="sdkjf"
        isDragging={false}
        selected={false}
        isConnectable={false}
      />
    </ReactFlowProvider>
  );
  expect(result.container).toMatchSnapshot();
});
