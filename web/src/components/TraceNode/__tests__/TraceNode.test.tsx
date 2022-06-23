import {render} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import TraceNode from '../TraceNode';
import {SemanticGroupNames} from '../../../constants/SemanticGroupNames.constants';

const data = {
  heading: '',
  name: '',
  primary: '',
  type: SemanticGroupNames.General,
  isAffected: false,
  isMatched: false,
};

test('TraceNode', () => {
  const {getByText} = render(
    <ReactFlowProvider>
      <TraceNode
        id="1234"
        type="span"
        data={data}
        selected={false}
        isConnectable={false}
        xPos={100}
        yPos={100}
        dragging={false}
        zIndex={1}
      />
    </ReactFlowProvider>
  );
  expect(getByText('General')).toBeTruthy();
});
