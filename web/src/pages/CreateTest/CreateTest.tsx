import Layout from 'components/Layout';
import CreateTestProvider from 'providers/CreateTest/CreateTest.provider';
import withAnalytics from '../../components/WithAnalytics/WithAnalytics';
import CreateTestContent from './CreateTestContent';

const CreateTestPage = () => {
  return (
    <Layout>
      <CreateTestProvider>
        <CreateTestContent />
      </CreateTestProvider>
    </Layout>
  );
};

export default withAnalytics(CreateTestPage, 'create-test');
