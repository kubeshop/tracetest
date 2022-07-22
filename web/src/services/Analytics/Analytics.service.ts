import {Categories} from '../../constants/Analytics.constants';

const {analyticsEnabled = 'false', serverId = ''} = window.ENV || {};
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
      label,
      category,
    });
  },
  page(name: string) {
    if (!isEnabled) return;
    analytics.page(name);
  },
  identify() {
    if (!isEnabled) return;
    analytics.identify({
      serverId,
    });
  },
});

export default AnalyticsService();
