import GA4React from 'ga-4-react';
import {Categories} from '../../constants/Analytics.constants';

const {analyticsEnabled = 'false', measurementId = ''} = window.ENV || {};

export const instance = new GA4React(measurementId);

export const isEnabled = analyticsEnabled === 'true';

type TAnalyticsService = {
  event<A>(category: Categories, action: A, label: string): void;
};

const AnalyticsService = (): TAnalyticsService => {
  const event = async <A>(category: Categories, action: A, label: string) => {
    if (!isEnabled) return;
    instance.event(String(action), label, category);
  };

  return {event};
};

export default AnalyticsService();
