import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useDataStoreConfig} from 'providers/DataStoreConfig/DataStoreConfig.provider';
import CreateTransactionProvider from 'providers/CreateTransaction';
import CreateTestProvider from 'providers/CreateTest';
import ConfigCTA from './ConfigCTA';
import Content from './Content';

const Home = () => {
  const {shouldDisplayConfigSetup, skipConfigSetup} = useDataStoreConfig();

  return (
    <Layout hasMenu>
      <CreateTransactionProvider>
        <CreateTestProvider>
          {shouldDisplayConfigSetup ? <ConfigCTA onSkip={skipConfigSetup} /> : <Content />}
        </CreateTestProvider>
      </CreateTransactionProvider>
    </Layout>
  );
};

export default withAnalytics(Home, 'home');
