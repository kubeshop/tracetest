import Layout from 'components/Layout';
import {useParams} from 'react-router-dom';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestProvider from 'providers/Test';
import TestContent from './TestContent';

const TestPage: React.FC = () => {
  const {testId = ''} = useParams();
  return (
    <Layout>
      <TestProvider testId={testId}>
        <TestContent />
      </TestProvider>
    </Layout>
  );
};

export default withAnalytics(TestPage, 'test-details');
