import posthog from 'posthog-js';
import {Categories} from 'constants/Analytics.constants';
import Env from 'utils/Env';

const isAnalyticsEnabled = () => Env.get('analyticsEnabled');
const appVersion = Env.get('appVersion');
const env = Env.get('env');
const serverID = Env.get('serverID');

const {analytics} = window;

type TAnalyticsService = {
  event<A>(category: Categories, action: A, label: string): void;
  page(page: string): void;
  identify(): void;
};

const AnalyticsService = (): TAnalyticsService => ({
  event<A>(category: Categories, action: A, label: string) {
    if (!isAnalyticsEnabled()) return;
    analytics.track(String(action), {
      serverID,
      appVersion,
      env,
      label,
      category,
    });
  },
  page(name: string) {
    if (!isAnalyticsEnabled()) return;
    analytics.page(name, {
      serverID,
      appVersion,
      env,
    });

    posthog.capture('$pageview');
  },
  identify() {
    if (!isAnalyticsEnabled()) return;
    analytics.identify({
      serverID,
      appVersion,
      env,
    });

    posthog.init('phc_Rg59ClPckoqa5p4onheukqHKJFPbTJkiNzECjIG4lMj', {
      api_host: 'https://app.posthog.com',
      loaded: ph => {
        ph.identify(serverID, {appVersion, env});
      },
    });
  },
});

export default AnalyticsService();
