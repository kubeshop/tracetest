import {render, waitFor} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {TestingModels} from '../../../utils/TestingModels';
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
          type={SupportedDiagrams.DAG}
          trace={TestingModels.trace}
          onSelectSpan={jest.fn()}
          selectedSpan={TestingModels.span}
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
          type={SupportedDiagrams.Timeline}
          trace={TestingModels.trace}
          onSelectSpan={jest.fn()}
          selectedSpan={TestingModels.span}
        />
      </ReactFlowProvider>
    </div>
  );

  await waitFor(() => getByText('Duration (ms)'));
});
