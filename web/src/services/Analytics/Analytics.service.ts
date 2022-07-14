import {Categories} from '../../constants/Analytics.constants';

const {analyticsEnabled = 'false'} = window.ENV || {};
const {analytics} = window;

export const isEnabled = analyticsEnabled === 'true';

type TAnalyticsService = {
  event<A>(category: Categories, action: A, label: string): void;
  page(page: string): void;
};

const AnalyticsService = (): TAnalyticsService => ({
  async event<A>(category: Categories, action: A, label: string) {
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
});

export default AnalyticsService();
