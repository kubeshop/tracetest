import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import CreateTestSuiteProvider from 'providers/CreateTestSuite';
import CreateTestProvider from 'providers/CreateTest';
import Content from './Content';

const Home = () => {
  const {isLoading, shouldDisplayConfigSetup, skipConfigSetup} = useSettingsValues();

  return (
    <CreateTestSuiteProvider>
      <CreateTestProvider>
        <Content
          isLoading={isLoading}
          shouldDisplayConfigSetup={shouldDisplayConfigSetup}
          skipConfigSetup={skipConfigSetup}
        />
      </CreateTestProvider>
    </CreateTestSuiteProvider>
  );
};

export default withAnalytics(Home, 'home');
