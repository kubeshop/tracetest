import {Categories} from 'constants/Analytics.constants';
import AnalyticsService from '../Analytics.service';

const eventMock = jest.fn();

Object.defineProperty(window, 'analytics', {
  event: eventMock,
} as any);

describe('AnalyticsService', () => {
  describe('event', () => {
    it('should not send an event if analyticsEnabled is false', async () => {
      await AnalyticsService.event(Categories.Home, 'test', 'test');

      expect(eventMock).not.toBeCalled();
    });
  });
});
