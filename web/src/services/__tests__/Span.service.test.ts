import {PseudoSelector} from '../../constants/Operator.constants';
import SpanMock from '../../models/__mocks__/Span.mock';
import SpanService from '../Span.service';

describe('SpanService', () => {
  describe('getSpanNodeInfo', () => {
    it('should return the span node info', () => {
      const span = SpanMock.model({
        attributes: {
          'db.system': 'mysql',
          'service.name': 'test',
          'tracetest.span.type': 'database',
        },
      });
      const info = SpanService.getSpanNodeInfo(span);

      expect(info.primary).toEqual('test');
      expect(info.heading).toEqual('mysql');
    });

    it('should handle empty attributes array', () => {
      const span = SpanMock.model({
        attributes: {},
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

      const selectorInfo = SpanService.getSelectorInformation(span);

      expect(selectorInfo.selectorList).toHaveLength(2);
      expect(selectorInfo.pseudoSelector).toEqual({selector: PseudoSelector.ALL});
    });
  });
});
