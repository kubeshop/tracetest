import Layout from 'components/Layout';
import {withTracker} from 'ga-4-react';
import EditTestModalProvider from 'components/EditTestModal/EditTestModal.provider';
import HomeContent from './HomeContent';

const Home = (): JSX.Element => {
  return (
    <Layout>
      <EditTestModalProvider>
        <HomeContent />
      </EditTestModalProvider>
    </Layout>
  );
};

export default withTracker(Home);
