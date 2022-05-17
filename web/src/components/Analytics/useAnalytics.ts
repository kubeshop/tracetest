import {useCallback, useContext} from 'react';
import {GA4ReactInterface} from 'ga-4-react/src/models/gtagModels';
import {Context} from './AnalyticsProvider';
import AnalyticsService from '../../services/Analytics/Analytics.service';
import {Categories} from '../../constants/Analytics.constants';

export type TAnalyticsService<A> = {
  isEnabled: boolean;
  instance: GA4ReactInterface;
  event(action: A, label: string): void;
};

const useAnalytics = <A>(category: Categories = Categories.Home): TAnalyticsService<A> => {
  const {instance, isEnabled} = useContext(Context);
  const event = useCallback((action, label) => AnalyticsService.event(category, action, label), [category]);

  return {isEnabled, event, instance};
};

export default useAnalytics;
