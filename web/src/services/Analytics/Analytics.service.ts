import {AnalyticsBrowser} from '@segment/analytics-next';
import {Categories} from 'constants/Analytics.constants';
import Env from 'utils/Env';

const isAnalyticsEnabled = () => Env.get('analyticsEnabled') && !Env.get('isTracetestDev');
const getServerID = () => Env.get('serverID');
const appVersion = Env.get('appVersion');
const env = Env.get('env');
const measurementId = Env.get('measurementId');

const getTraits = () => ({
  serverID: getServerID(),
  ...(appVersion && {appVersion}),
  ...(env && {env}),
});

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
      ...getTraits(),
      label,
      category,
    });
  },
  page(name: string) {
    if (!isAnalyticsEnabled()) return;
    analytics.page(name, getTraits());
  },
  identify() {
    if (!isAnalyticsEnabled()) return;
    analytics.identify(getServerID(), getTraits());
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
      Env.set('segmentLoaded', false);
    }
  },
});

export default AnalyticsService();
