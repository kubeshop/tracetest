import {Categories} from 'constants/Analytics.constants';
import Env from 'utils/Env';

const analyticsEnabled = Env.get('analyticsEnabled');
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
    if (!analyticsEnabled) return;
    analytics.track(String(action), {
      serverID,
      appVersion,
      env,
      label,
      category,
    });
  },
  page(name: string) {
    if (!analyticsEnabled) return;
    analytics.page(name, {
      serverID,
      appVersion,
      env,
    });
  },
  identify() {
    if (!analyticsEnabled) return;
    analytics.identify({
      serverID,
      appVersion,
      env,
    });
  },
});

export default AnalyticsService();
