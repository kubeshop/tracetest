import Layout from 'components/Layout';
import EditTestModalProvider from '../../components/EditTestModal/EditTestModal.provider';
import TestContent from './TestContent';
import withAnalytics from '../../components/WithAnalytics/WithAnalytics';

const TestPage: React.FC = () => {
  return (
    <Layout>
      <EditTestModalProvider>
        <TestContent />
      </EditTestModalProvider>
    </Layout>
  );
};

export default withAnalytics(TestPage, 'test-details');
