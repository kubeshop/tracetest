import {useParams} from 'react-router-dom';

import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestSuiteProvider from 'providers/TestSuite';
import Content from './Content';

const TestSuite = () => {
  const {testSuiteId = ''} = useParams();

  return (
    <Layout hasMenu>
      <TestSuiteProvider testSuiteId={testSuiteId}>
        <Content />
      </TestSuiteProvider>
    </Layout>
  );
};

export default withAnalytics(TestSuite, 'testsuite-details');
