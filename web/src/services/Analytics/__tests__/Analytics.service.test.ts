import {Categories} from '../../../constants/Analytics.constants';
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
    it('should handle sending an event', async () => {
      expect.assertions(2);

      await AnalyticsService.event(Categories.Home, 'test', 'test');

      expect(instance.event).toBeCalledTimes(1);
      expect(instance.event).toBeCalledWith('test', 'test', Categories.Home);
    });
  });
});
