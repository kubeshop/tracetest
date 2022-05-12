import {GA4React} from 'ga-4-react';
import {useCallback, useContext} from 'react';
import {Categories} from '../../constants/Analytics.constants';
import AnalyticsService from '../../services/Analytics/Analytics.service';
import {Context} from './AnalyticsProvider';

export type TAnalyticsService<A> = {
  isEnabled: boolean;
  instance: GA4React;
  event(action: A, label: string): void;
};

const useAnalytics = <A>(category: Categories = Categories.Home): TAnalyticsService<A> => {
  const {instance, isEnabled} = useContext(Context);
  const event = useCallback((action, label) => AnalyticsService.event(category, action, label), [category]);

  return {isEnabled, event, instance};
};

export default useAnalytics;
