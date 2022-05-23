import {withTracker} from 'ga-4-react';
import {useParams} from 'react-router-dom';
import {ReactFlowProvider} from 'react-flow-renderer';
import Layout from 'components/Layout';
import AssertionFormProvider from '../../components/AssertionForm/AssertionFormProvider';
import TestRunProvider from '../../providers/TestRun';
import TraceContent from './TraceContent';

const TracePage = () => {
  const {testId = '', runId = ''} = useParams();

  return (
    <ReactFlowProvider>
      <AssertionFormProvider testId={testId}>
        <TestRunProvider testId={testId} runId={runId}>
          <Layout>
            <TraceContent />
          </Layout>
        </TestRunProvider>
      </AssertionFormProvider>
    </ReactFlowProvider>
  );
};

export default withTracker(TracePage);
