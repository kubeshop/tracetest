import {render, waitFor} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import Trace from '../Trace';

test('Trace', async () => {
  const {getByText} = render(
    <ReactFlowProvider>
      <ReduxWrapperProvider>
        <div style={{width: 600, height: 600}}>
          <Trace
            minHeight="300px"
            run={TestRunMock.model()}
            test={TestMock.model()}
            visiblePortion={100}
            displayError={false}
          />
        </div>
      </ReduxWrapperProvider>
    </ReactFlowProvider>
  );

  await waitFor(() => getByText('HTTP'));
});
