import SpanMock from '../../models/__mocks__/Span.mock';
import SpanService from '../Span.service';

describe('SpanService', () => {
  describe('getSelectedSpanListAttributes', () => {
    it('should return the selected span list attributes', () => {
      const span = SpanMock.model();

      const {intersectedList, differenceList} = SpanService.getSelectedSpanListAttributes(span, [span]);

      expect(span.attributeList).toHaveLength(intersectedList.length);
      expect(differenceList).toHaveLength(0);
    });

    it('should return the selected span list attributes with different sizes', () => {
      const span = SpanMock.model();
      const spanList = [
        SpanMock.model(),
        SpanMock.model({
          attributes: {
            'db.system': 'mysql',
            'service.name': 'mock',
          },
        }),
      ];

      const {differenceList, intersectedList} = SpanService.getSelectedSpanListAttributes(span, spanList);

      expect(intersectedList.length).toBeGreaterThan(0);
      expect(differenceList.length).toBeGreaterThan(0);
    });
  });

  describe('getSelectorInformation', () => {
    it('should return the spanList', () => {
      const span = SpanMock.model();

      const selectorInfo = SpanService.searchSpanList([span], 'tracetest');

      expect(selectorInfo).toHaveLength(1);
    });
  });
  describe('getSelectorInformation', () => {
    it('should return the selector information', () => {
      const span = SpanMock.model();

      const selector = SpanService.getSelectorInformation(span);

      expect(typeof selector === 'string').toBeTruthy();
    });
  });
});
