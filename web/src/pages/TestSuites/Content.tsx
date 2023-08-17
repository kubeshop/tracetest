import ConfigCTA from '../Home/ConfigCTA';
import TestSuites from './TestSuitesList';

interface IProps {
  isLoading: boolean;
  shouldDisplayConfigSetup: boolean;
  skipConfigSetup(): void;
}

const Content = ({isLoading, shouldDisplayConfigSetup, skipConfigSetup}: IProps) => {
  if (isLoading) return null;

  return shouldDisplayConfigSetup ? <ConfigCTA onSkip={skipConfigSetup} /> : <TestSuites />;
};

export default Content;
