import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useDataStoreConfig} from 'providers/DataStoreConfig/DataStoreConfig.provider';
import CreateTransactionProvider from 'providers/CreateTransaction';
import CreateTestProvider from 'providers/CreateTest';
import Content from './Content';

const Home = () => {
  const {isLoading, shouldDisplayConfigSetup, skipConfigSetup} = useDataStoreConfig();

  return (
    <Layout hasMenu>
      <CreateTransactionProvider>
        <CreateTestProvider>
          <Content
            isLoading={isLoading}
            shouldDisplayConfigSetup={shouldDisplayConfigSetup}
            skipConfigSetup={skipConfigSetup}
          />
        </CreateTestProvider>
      </CreateTransactionProvider>
    </Layout>
  );
};

export default withAnalytics(Home, 'home');
