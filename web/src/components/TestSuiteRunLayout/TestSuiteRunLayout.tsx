import {useParams} from 'react-router-dom';
import Layout from 'components/Layout';
import TestSuiteHeader from 'components/TestSuiteHeader';
import TestSuiteRunProvider from 'providers/TestSuiteRun/TestSuite.provider';

interface IProps {
  children: React.ReactNode;
}

const TestSuiteRunLayout = ({children}: IProps) => {
  const {testSuiteId = '', runId = ''} = useParams();

  return (
    <Layout>
      <TestSuiteRunProvider testSuiteId={testSuiteId} runId={runId}>
        <TestSuiteHeader />
        {children}
      </TestSuiteRunProvider>
    </Layout>
  );
};

export default TestSuiteRunLayout;
