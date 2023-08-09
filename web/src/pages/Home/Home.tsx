import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import CreateTestSuiteProvider from 'providers/CreateTestSuite';
import CreateTestProvider from 'providers/CreateTest';
import Content from './Content';

const Home = () => {
  const {isLoading, shouldDisplayConfigSetup, skipConfigSetup} = useSettingsValues();

  return (
    <Layout hasMenu>
      <CreateTestSuiteProvider>
        <CreateTestProvider>
          <Content
            isLoading={isLoading}
            shouldDisplayConfigSetup={shouldDisplayConfigSetup}
            skipConfigSetup={skipConfigSetup}
          />
        </CreateTestProvider>
      </CreateTestSuiteProvider>
    </Layout>
  );
};

export default withAnalytics(Home, 'home');
