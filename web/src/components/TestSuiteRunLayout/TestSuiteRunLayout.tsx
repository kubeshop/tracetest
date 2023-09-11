import {useParams} from 'react-router-dom';
import TestSuiteHeader from 'components/TestSuiteHeader';
import TestSuiteRunProvider from 'providers/TestSuiteRun/TestSuite.provider';

interface IProps {
  children: React.ReactNode;
}

const TestSuiteRunLayout = ({children}: IProps) => {
  const {testSuiteId = '', runId = 0} = useParams();

  return (
    <TestSuiteRunProvider testSuiteId={testSuiteId} runId={Number(runId)}>
      <TestSuiteHeader />
      {children}
    </TestSuiteRunProvider>
  );
};

export default TestSuiteRunLayout;
