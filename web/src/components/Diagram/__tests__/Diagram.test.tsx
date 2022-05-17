import {render} from '@testing-library/react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {TestingModels} from '../../../utils/TestingModels';
import AnalyticsProvider from '../../Analytics';
import Diagram, {SupportedDiagrams} from '../Diagram';

test('Diagram DAG', () => {
  const result = render(
    <AnalyticsProvider>
      <ReactFlowProvider>
        <Diagram
          type={SupportedDiagrams.DAG}
          trace={TestingModels.trace}
          onSelectSpan={jest.fn()}
          selectedSpan={TestingModels.span}
        />
      </ReactFlowProvider>
    </AnalyticsProvider>
  );
  expect(result.container).toMatchSnapshot();
});
test('Diagram Timeline', () => {
  const result = render(
    <AnalyticsProvider>
      <ReactFlowProvider>
        <Diagram
          type={SupportedDiagrams.Timeline}
          trace={TestingModels.trace}
          onSelectSpan={jest.fn()}
          selectedSpan={TestingModels.span}
        />
      </ReactFlowProvider>
    </AnalyticsProvider>
  );
  expect(result.container).toMatchSnapshot();
});
