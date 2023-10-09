import {usePostHog} from 'posthog-js/react';
import CaptureProvider from 'providers/Capture';
import {useCallback, useMemo} from 'react';
import Env from 'utils/Env';

const appVersion = Env.get('appVersion');
const env = Env.get('env');
const serverID = Env.get('serverID');
const isAnalyticsEnabled = () => Env.get('analyticsEnabled') && !Env.get('isTracetestDev');

interface IProps {
  children: React.ReactNode;
}

const Content = ({children}: IProps) => {
  const posthog = usePostHog();

  const identify = useCallback(() => {
    console.log('identify', 'isAnalyticsEnabled', isAnalyticsEnabled());
    if (!isAnalyticsEnabled()) {
      return;
    }

    posthog?.opt_in_capturing();

    posthog?.identify(serverID, {
      appVersion,
      env,
    });
  }, [posthog]);

  const pageView = useCallback(() => {
    posthog?.capture('$pageview');
  }, [posthog]);

  const captureProviderValue = useMemo(() => ({identify, pageView}), [identify, pageView]);

  return <CaptureProvider value={captureProviderValue}>{children}</CaptureProvider>;
};

export default Content;
