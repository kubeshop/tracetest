import {AnalyticsBrowser} from '@segment/analytics-next';
import posthog from 'posthog-js';
import {Categories} from 'constants/Analytics.constants';
import Env from 'utils/Env';

const isAnalyticsEnabled = () => Env.get('analyticsEnabled') && !Env.get('isTracetestDev');
const appVersion = Env.get('appVersion');
const env = Env.get('env');
const serverID = Env.get('serverID');
const measurementId = Env.get('measurementId');
const posthogKey = Env.get('posthogKey');

export const analytics = new AnalyticsBrowser();

type TAnalyticsService = {
  event<A>(category: Categories, action: A, label: string): void;
  page(page: string): void;
  identify(): void;
  load(): void;
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

    posthog.init(posthogKey, {
      api_host: 'https://app.posthog.com',
      loaded: ph => {
        ph.identify(serverID, {appVersion, env});
      },
    });

    if (posthog.has_opted_out_capturing()) {
      posthog.opt_in_capturing();
    }
  },
  load() {
    const isSegmentLoaded = Env.get('segmentLoaded');

    if (isAnalyticsEnabled() && !isSegmentLoaded) {
      analytics.load({writeKey: measurementId});
      Env.set('segmentLoaded', true);
      return;
    }

    if (!isAnalyticsEnabled() && isSegmentLoaded) {
      analytics.reset();
      posthog.persistence && posthog.opt_out_capturing();
      Env.set('segmentLoaded', false);
    }
  },
});

export default AnalyticsService();
