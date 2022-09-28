import {useParams} from 'react-router-dom';

import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestProvider from 'providers/Test';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import TestContent from './TestContent';

const Test = () => {
  const {testId = ''} = useParams();

  return (
    <Layout hasMenu={ExperimentalFeature.isEnabled('transactions')}>
      <TestProvider testId={testId}>
        <TestContent />
      </TestProvider>
    </Layout>
  );
};

export default withAnalytics(Test, 'test-details');
