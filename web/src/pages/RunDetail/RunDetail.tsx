import {useParams} from 'react-router-dom';
import {ReactFlowProvider} from 'react-flow-renderer';

import TestSpecFormProvider from 'components/TestSpecForm/TestSpecForm.provider';
import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import GuidedTourProvider from 'providers/GuidedTour/GuidedTour.provider';
import SpanProvider from 'providers/Span';
import TestDefinitionProvider from 'providers/TestDefinition';
import TestRunProvider from 'providers/TestRun';
import Content from './Content';

const RunDetail = () => {
  const {testId = '', runId = ''} = useParams();

  return (
    <GuidedTourProvider>
      <ReactFlowProvider>
        <TestRunProvider testId={testId} runId={runId}>
          <TestDefinitionProvider testId={testId} runId={runId}>
            <TestSpecFormProvider testId={testId}>
              <SpanProvider>
                <Layout>
                  <Content />
                </Layout>
              </SpanProvider>
            </TestSpecFormProvider>
          </TestDefinitionProvider>
        </TestRunProvider>
      </ReactFlowProvider>
    </GuidedTourProvider>
  );
};

export default withAnalytics(RunDetail, 'run-detail');
