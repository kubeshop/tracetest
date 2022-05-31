import SpanAttributeService from '../SpanAttribute.service';
import {SemanticGroupNames} from '../../constants/SemanticGroupNames.constants';
import {SectionNames} from '../../constants/Span.constants';
import {Attributes} from '../../constants/SpanAttribute.constants';

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
          section: SectionNames.operation,
          attributeList: [],
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
          section: SectionNames.producer,
          attributeList: [],
        },
        {
          section: SectionNames.consumer,
          attributeList: [],
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
});
