import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import CreateTestSuiteProvider from 'providers/CreateTestSuite';
import Content from './Content';

const TestSuites = () => {
  const {isLoading, shouldDisplayConfigSetup, skipConfigSetup} = useSettingsValues();

  return (
    <CreateTestSuiteProvider>
      <Content
        isLoading={isLoading}
        shouldDisplayConfigSetup={shouldDisplayConfigSetup}
        skipConfigSetup={skipConfigSetup}
      />
    </CreateTestSuiteProvider>
  );
};

export default withAnalytics(TestSuites, 'test-suite');
