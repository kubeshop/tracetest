import Layout from 'components/Layout';
import TestContent from './TestContent';
import withAnalytics from '../../components/WithAnalytics/WithAnalytics';

const TestPage: React.FC = () => {
  return (
    <Layout>
      <TestContent />
    </Layout>
  );
};

export default withAnalytics(TestPage, 'test-details');
