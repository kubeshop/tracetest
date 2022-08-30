import {render, waitFor} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import Run from '../Run';

test.skip('Run', async () => {
  const {getByText} = render(
    <ReactFlowProvider>
      <ReduxWrapperProvider>
        <div style={{width: 600, height: 600}}>
          <Run run={TestRunMock.model()} />
        </div>
      </ReduxWrapperProvider>
    </ReactFlowProvider>
  );

  await waitFor(() => getByText('HTTP'));
});
