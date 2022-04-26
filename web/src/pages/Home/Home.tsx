import Layout from 'components/Layout';
import {withTracker} from 'ga-4-react';
import HomeContent from './HomeContent';

const Home = (): JSX.Element => {
  return (
    <Layout>
      <HomeContent />
    </Layout>
  );
};

export default withTracker(Home);
