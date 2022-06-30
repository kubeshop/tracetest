import SpanAttributeService from '../SpanAttribute.service';
import {SemanticGroupNames} from '../../constants/SemanticGroupNames.constants';
import {SectionNames} from '../../constants/Span.constants';
import {Attributes} from '../../constants/SpanAttribute.constants';
import {TSpanFlatAttribute} from '../../types/Span.types';

describe('SpanAttributeService', () => {
  describe('getFilteredSpanAttributeList', () => {
    it('should return the sections for http type', () => {
      expect(SpanAttributeService.getSpanAttributeSectionsList([], SemanticGroupNames.Http)).toEqual([
        {
          section: SectionNames.request,
          attributeList: [],
        },
        {
          section: SectionNames.response,
          attributeList: [],
        },
        {
          section: SectionNames.custom,
          attributeList: [],
        },
        {
          section: SectionNames.all,
          attributeList: [],
        },
      ]);
    });

    it('should return the sections for database type with values', () => {
      const attribute = {
        key: Attributes.DB_SYSTEM,
        value: 'mysql',
      };

      expect(SpanAttributeService.getSpanAttributeSectionsList([attribute], SemanticGroupNames.Database)).toEqual([
        {
          section: SectionNames.metadata,
          attributeList: [attribute],
        },
        {
          section: SectionNames.custom,
          attributeList: [],
        },
        {
          section: SectionNames.all,
          attributeList: [attribute],
        },
      ]);
    });

    it('should return the sections for messaging type with values', () => {
      const attribute = {
        key: Attributes.MESSAGING_SYSTEM,
        value: 'kafka',
      };

      const customAttribute = {
        key: 'messaging.payload',
        value: '{}',
      };

      expect(
        SpanAttributeService.getSpanAttributeSectionsList([attribute, customAttribute], SemanticGroupNames.Messaging)
      ).toEqual([
        {
          section: SectionNames.metadata,
          attributeList: [attribute],
        },
        {
          section: SectionNames.custom,
          attributeList: [customAttribute],
        },
        {
          section: SectionNames.all,
          attributeList: [attribute, customAttribute],
        },
      ]);
    });
  });

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
