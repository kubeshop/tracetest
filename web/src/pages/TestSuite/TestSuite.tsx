import {useParams} from 'react-router-dom';

import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestSuiteProvider from 'providers/TestSuite';
import Content from './Content';

const TestSuite = () => {
  const {testSuiteId = ''} = useParams();

  return (
    <TestSuiteProvider testSuiteId={testSuiteId}>
      <Content />
    </TestSuiteProvider>
  );
};

export default withAnalytics(TestSuite, 'testsuite-details');
