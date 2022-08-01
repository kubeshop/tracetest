import {Categories} from '../../constants/Analytics.constants';

const {analyticsEnabled = 'false', serverID = '', appVersion = '', env = ''} = window.ENV || {};
const {analytics} = window;

export const isEnabled = analyticsEnabled === 'true';

type TAnalyticsService = {
  event<A>(category: Categories, action: A, label: string): void;
  page(page: string): void;
  identify(): void;
};

const AnalyticsService = (): TAnalyticsService => ({
  event<A>(category: Categories, action: A, label: string) {
    if (!isEnabled) return;
    analytics.track(String(action), {
      serverID,
      appVersion,
      env,
      label,
      category,
    });
  },
  page(name: string) {
    if (!isEnabled) return;
    analytics.page(name, {
      serverID,
      appVersion,
      env,
    });
  },
  identify() {
    if (!isEnabled) return;
    analytics.identify({
      serverID,
      appVersion,
      env,
    });
  },
});

export default AnalyticsService();
