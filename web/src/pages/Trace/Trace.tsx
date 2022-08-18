import {useParams} from 'react-router-dom';
import {ReactFlowProvider} from 'react-flow-renderer';
import Layout from 'components/Layout';
import AssertionFormProvider from 'components/AssertionForm/AssertionForm.provider';
import TestRunProvider from 'providers/TestRun';
import TestDefinitionProvider from 'providers/TestDefinition';
import GuidedTourProvider from 'providers/GuidedTour/GuidedTour.provider';
import SpanProvider from 'providers/Span';
import TraceContent from './TraceContent';
import withAnalytics from '../../components/WithAnalytics/WithAnalytics';

const TracePage = () => {
  const {testId = '', runId = ''} = useParams();

  return (
    <Layout>
      <GuidedTourProvider>
        <ReactFlowProvider>
          <TestRunProvider testId={testId} runId={runId}>
            <TestDefinitionProvider testId={testId} runId={runId}>
              <AssertionFormProvider testId={testId}>
                <SpanProvider>
                  <TraceContent />
                </SpanProvider>
              </AssertionFormProvider>
            </TestDefinitionProvider>
          </TestRunProvider>
        </ReactFlowProvider>
      </GuidedTourProvider>
    </Layout>
  );
};

export default withAnalytics(TracePage, 'trace');
