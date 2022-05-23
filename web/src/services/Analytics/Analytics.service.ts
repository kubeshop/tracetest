import GA4React from 'ga-4-react';
import {Categories} from '../../constants/Analytics.constants';

const {analyticsEnabled = 'true', measurementId = ''} = window.ENV || {};

export const instance = new GA4React(measurementId);

export const isEnabled = analyticsEnabled === 'true';

type TAnalyticsService = {
  event<A>(category: Categories, action: A, label: string): void;
};

const initializePromise = instance.initialize();

const AnalyticsService = (): TAnalyticsService => {
  const event = async <A>(category: Categories, action: A, label: string) => {
    await initializePromise;
    instance.event(String(action), label, category);
  };

  return {event};
};

export default AnalyticsService();
