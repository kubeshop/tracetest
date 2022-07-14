import Layout from 'components/Layout';
import EditTestModalProvider from 'components/EditTestModal/EditTestModal.provider';
import HomeContent from './HomeContent';
import withAnalytics from '../../components/WithAnalytics/WithAnalytics';

const Home = (): JSX.Element => {
  return (
    <Layout>
      <EditTestModalProvider>
        <HomeContent />
      </EditTestModalProvider>
    </Layout>
  );
};

export default withAnalytics(Home, 'home');
