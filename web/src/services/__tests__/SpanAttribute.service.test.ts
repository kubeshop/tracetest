import SpanAttributeService from '../SpanAttribute.service';
import {Attributes} from '../../constants/SpanAttribute.constants';
import {TSpanFlatAttribute} from '../../types/Span.types';

describe('SpanAttributeService', () => {
  describe('getFilteredSelectorAttributeList', () => {
    it('should return the filtered list of attributes', () => {
      const attributeList: TSpanFlatAttribute[] = [
        {
          key: Attributes.MESSAGING_SYSTEM,
          value: 'kafka',
        },
        {
          key: Attributes.DB_NAME,
          value: 'pokeshop',
        },
        {
          key: Attributes.DB_STATEMENT,
          value: 'SELECT * FROM users',
        },
        {
          key: Attributes.TRACETEST_RESPONSE_BODY,
          value: '{"id": 1}',
        },
      ];

      expect(SpanAttributeService.getFilteredSelectorAttributeList(attributeList, [])).toEqual([
        {
          key: Attributes.DB_NAME,
          value: 'pokeshop',
        },
        {
          key: Attributes.MESSAGING_SYSTEM,
          value: 'kafka',
        },
      ]);
    });
  });
});
