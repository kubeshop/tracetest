import {useParams} from 'react-router-dom';

import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestProvider from 'providers/Test';
import Content from './Content';

const Test = () => {
  const {testId = ''} = useParams();

  return (
    <Layout hasMenu>
      <TestProvider testId={testId}>
        <Content />
      </TestProvider>
    </Layout>
  );
};

export default withAnalytics(Test, 'test-details');
