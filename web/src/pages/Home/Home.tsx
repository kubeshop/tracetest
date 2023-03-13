import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import CreateTransactionProvider from 'providers/CreateTransaction';
import CreateTestProvider from 'providers/CreateTest';
import Content from './Content';

const Home = () => {
  const {isLoading, shouldDisplayConfigSetup, skipConfigSetup} = useSettingsValues();

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
