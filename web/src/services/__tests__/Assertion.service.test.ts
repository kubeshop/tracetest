import faker from '@faker-js/faker';
import AssertionResultMock from '../../models/__mocks__/AssertionResult.mock';
import AssertionSpanResultMock from '../../models/__mocks__/AssertionSpanResult.mock';
import AssertionService from '../Assertion.service';

describe('AssertionService', () => {
  describe('extractExpectedString', () => {
    it('should return quoted string', () => {
      expect(AssertionService.extractExpectedString('some text')).toEqual(AssertionService.quotedString('some text'));
    });
    it('should return not quoted string', () => {
      const {extractExpectedString} = AssertionService;
      expect(extractExpectedString(undefined)).toEqual(undefined);
      expect(extractExpectedString('')).toEqual('');
      expect(extractExpectedString('300')).toEqual('300');
      expect(extractExpectedString('300.3')).toEqual('300.3');
      expect(extractExpectedString('300ms')).toEqual('300ms');
      expect(extractExpectedString('300.3ms')).toEqual('300.3ms');
      expect(extractExpectedString('tracetest.span.duration')).toEqual('tracetest.span.duration');
      expect(extractExpectedString('tracetest.span.duration + 30ms')).toEqual('tracetest.span.duration + 30ms');
    });
  });
  describe('getSpanIds', () => {
    it('should return a list of unique span ids', () => {
      const {getSpanIds} = AssertionService;
      const id1 = faker.datatype.uuid();
      const id2 = faker.datatype.uuid();
      const id3 = faker.datatype.uuid();
      const assertionResult = AssertionResultMock.model({
        spanResults: [
          AssertionSpanResultMock.raw({spanId: id1}),
          AssertionSpanResultMock.raw({spanId: id1}),
          AssertionSpanResultMock.raw({spanId: id2}),
          AssertionSpanResultMock.raw({spanId: id3}),
        ],
      });
      const result = getSpanIds([assertionResult]);

      expect(result).toEqual([id1, id2, id3]);
    });
  });
});
