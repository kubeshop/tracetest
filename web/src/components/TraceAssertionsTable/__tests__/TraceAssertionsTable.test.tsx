import {fireEvent, render} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import {TestingModels, TestingObj} from '../../../utils/TestingModels';
import TraceAssertionsTable from '../TraceAssertionsTable';

test('TraceAssertionsTable', () => {
  const result = render(
    <ReactFlowProvider>
      <ReduxWrapperProvider>
        <TraceAssertionsTable assertionResult={TestingObj.assertionResult} onSpanSelected={jest.fn()} />
      </ReduxWrapperProvider>
    </ReactFlowProvider>
  );
  expect(result.container).toMatchSnapshot();
});

test('TraceAssertionsTable click row', () => {
  const result = render(
    <ReactFlowProvider>
      <ReduxWrapperProvider>
        <TraceAssertionsTable assertionResult={TestingObj.assertionResult} onSpanSelected={jest.fn()} />
      </ReduxWrapperProvider>
    </ReactFlowProvider>
  );

  const elementsByClassName = result.container.getElementsByClassName('ant-table-row');
  expect(result.container).toMatchSnapshot();
  fireEvent(elementsByClassName.item(0)!, TestingModels.mouseEvent);
});
