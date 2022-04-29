import {useContext, useMemo} from 'react';
import {GA4ReactInterface} from 'ga-4-react/src/models/gtagModels';
import {Context} from './AnalyticsProvider';
import AnalyticsService, {Categories} from '../../services/Analytics/Analytics.service';

export type TAnalyticsService<A> = {
  isEnabled: boolean;
  instance: GA4ReactInterface;
  event(action: A, label: string): void;
};

const useAnalytics = <A>(category: Categories = Categories.Home): TAnalyticsService<A> => {
  const {instance, isEnabled} = useContext(Context);
  const {event} = useMemo(() => AnalyticsService(category), [category]);

  return {isEnabled, event, instance};
};

export default useAnalytics;
