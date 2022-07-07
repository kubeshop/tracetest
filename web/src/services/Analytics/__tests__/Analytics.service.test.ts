import {Categories} from 'constants/Analytics.constants';
import AnalyticsService, {instance} from '../Analytics.service';

jest.mock('ga-4-react', () => {
  return jest.fn(() => {
    return {
      initialize: jest.fn(() => Promise.resolve()),
      event: jest.fn(),
    };
  });
});

describe('AnalyticsService', () => {
  describe('event', () => {
    it('should not send an event if analyticsEnabled is false', async () => {
      await AnalyticsService.event(Categories.Home, 'test', 'test');

      expect(instance.event).not.toBeCalled();
    });
  });
});
