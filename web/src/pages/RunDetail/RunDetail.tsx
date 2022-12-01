import {useParams} from 'react-router-dom';

import Layout from 'components/Layout';
import TestSpecFormProvider from 'components/TestSpecForm/TestSpecForm.provider';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import GuidedTourProvider from 'providers/GuidedTour';
import TestProvider from 'providers/Test';
import TestOutputProvider from 'providers/TestOutput';
import SpanProvider from 'providers/Span';
import TestSpecsProvider from 'providers/TestSpecs';
import TestRunProvider from 'providers/TestRun';
import Content from './Content';

const RunDetail = () => {
  const {testId = '', runId = ''} = useParams();

  return (
    <GuidedTourProvider>
      <Layout>
        <TestRunProvider testId={testId} runId={runId}>
          <TestProvider testId={testId}>
            <TestSpecsProvider testId={testId} runId={runId}>
              <TestSpecFormProvider testId={testId}>
                <SpanProvider>
                  <TestOutputProvider testId={testId} runId={runId}>
                    <Content />
                  </TestOutputProvider>
                </SpanProvider>
              </TestSpecFormProvider>
            </TestSpecsProvider>
          </TestProvider>
        </TestRunProvider>
      </Layout>
    </GuidedTourProvider>
  );
};

export default withAnalytics(RunDetail, 'run-detail');
