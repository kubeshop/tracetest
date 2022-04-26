import {useCallback, useContext} from 'react';
import {GA4ReactInterface} from 'ga-4-react/src/models/gtagModels';
import {Context} from './AnalyticsProvider';

export enum Categories {
  Home = 'home',
  Test = 'test',
  Trace = 'trace',
  TraceDetail = 'trace-detail',
  TestResults = 'test-results',
  SpanDetail = 'span-detail',
  Assertion = 'assertion',
}

export enum Labels {
  Button = 'button',
  Link = 'link',
  Modal = 'modal',
  Table = 'table',
  Form = 'form',
}

export type TAnalyticsService<A> = {
  isEnabled: boolean;
  instance: GA4ReactInterface;
  event(action: A, label: string): void;
};

const useAnalytics = <A>(category: Categories = Categories.Home): TAnalyticsService<A> => {
  const {instance, isEnabled} = useContext(Context);

  const event = useCallback(
    (action: A, label: string) => {
      instance.event(String(action), label, category);
    },
    [category, instance]
  );

  return {isEnabled, event, instance};
};

export default useAnalytics;
