import GA4React from 'ga-4-react';

const {analyticsEnabled = 'true', measurementId = 'G-ZP277L2M37'} = window.ENV || {};

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

export const instance = new GA4React(measurementId);

export const isEnabled = analyticsEnabled === 'true';

type TAnalyticsService<A> = {
  event(action: A, label: string): void;
};

const AnalyticsService = <A>(category: Categories): TAnalyticsService<A> => {
  const event = (action: A, label: string) => {
    instance.event(String(action), label, category);
  };

  return {event};
};

export default AnalyticsService;
