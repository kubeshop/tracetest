import SpanMock from '../../models/__mocks__/Span.mock';
import SpanService from '../Span.service';

describe('SpanService', () => {
  describe('getSpanNodeInfo', () => {
    it('should return the span node info', () => {
      const span = SpanMock.model({
        attributes: [
          {
            key: 'db.system',
            value: {
              stringValue: 'mysql',
              kvlistValue: {values: []},
            },
          },
          {
            key: 'service.name',
            value: {
              stringValue: 'test',
              kvlistValue: {values: []},
            },
          },
        ],
      });
      const info = SpanService.getSpanNodeInfo(span);

      expect(info.primary).toEqual('test');
      expect(info.heading).toEqual('mysql');
    });

    it('should handle empty attributes array', () => {
      const span = SpanMock.model({
        attributes: [],
      });
      const info = SpanService.getSpanNodeInfo(span);

      expect(info.primary).toEqual('');
      expect(info.heading).toEqual('');
    });
  });

  describe('getSelectedSpanListAttributes', () => {
    it('should return the selected span list attributes', () => {
      const span = SpanMock.model();

      const {intersectedList, differenceList} = SpanService.getSelectedSpanListAttributes(span, [span]);

      expect(span.attributeList).toHaveLength(intersectedList.length);
      expect(differenceList).toHaveLength(0);
    });

    it('should return the selected span list attributes with different sizes', () => {
      const span = SpanMock.model();
      const spanList = [SpanMock.model(), SpanMock.model()];

      const {differenceList, intersectedList} = SpanService.getSelectedSpanListAttributes(span, spanList);

      expect(intersectedList.length).toBeGreaterThan(0);
      expect(differenceList.length).toBeGreaterThan(0);
    });
  });
});
