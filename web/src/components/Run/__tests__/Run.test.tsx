import {render, waitFor} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import Run from '../Run';

test.skip('Run', async () => {
  const {getByText} = render(
    <ReactFlowProvider>
      <ReduxWrapperProvider>
        <div style={{width: 600, height: 600}}>
          <Run displayError={false} run={TestRunMock.model()} test={TestMock.model()} />
        </div>
      </ReduxWrapperProvider>
    </ReactFlowProvider>
  );

  await waitFor(() => getByText('HTTP'));
});
