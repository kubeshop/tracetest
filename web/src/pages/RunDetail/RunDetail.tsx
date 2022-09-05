import {useParams} from 'react-router-dom';

import Layout from 'components/Layout';
import TestSpecFormProvider from 'components/TestSpecForm/TestSpecForm.provider';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import GuidedTourProvider from 'providers/GuidedTour/GuidedTour.provider';
import SpanProvider from 'providers/Span';
import TestDefinitionProvider from 'providers/TestDefinition';
import TestRunProvider from 'providers/TestRun';
import TestProvider from 'providers/Test';
import Content from './Content';

const RunDetail = () => {
  const {testId = '', runId = ''} = useParams();

  return (
    <GuidedTourProvider>
      <Layout>
        <TestRunProvider testId={testId} runId={runId}>
          <TestDefinitionProvider testId={testId} runId={runId}>
            <TestProvider testId={testId}>
              <TestSpecFormProvider testId={testId}>
                <SpanProvider>
                  <Content />
                </SpanProvider>
              </TestSpecFormProvider>
            </TestProvider>
          </TestDefinitionProvider>
        </TestRunProvider>
      </Layout>
    </GuidedTourProvider>
  );
};

export default withAnalytics(RunDetail, 'run-detail');
