import {render, waitFor} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import {TestingModels} from '../../../utils/TestingModels';
import Trace from '../Trace';

test('Trace', async () => {
  const {getByText} = render(
    <ReactFlowProvider>
      <ReduxWrapperProvider>
        <div style={{width: 600, height: 600}}>
          <Trace
            minHeight="300px"
            testResultDetails={TestingModels.testRunResult}
            test={TestingModels.test}
            visiblePortion={100}
            displayError={false}
          />
        </div>
      </ReduxWrapperProvider>
    </ReactFlowProvider>
  );

  await waitFor(() => getByText('HTTP'));
});
