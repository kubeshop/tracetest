import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import HomeContent from './HomeContent';
import LayoutV2 from './Layout';

const {featuresEnabled = '[]'} = window.ENV || {};
const parsedFeaturesEnabled = JSON.parse(featuresEnabled);
const isTransactionsEnabled = parsedFeaturesEnabled?.includes('transactions');

const Home = (): JSX.Element => {
  return isTransactionsEnabled ? (
    <LayoutV2>
      <HomeContent />
    </LayoutV2>
  ) : (
    <Layout>
      <HomeContent />
    </Layout>
  );
};

export default withAnalytics(Home, 'home');
