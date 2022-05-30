import faker from '@faker-js/faker';
import AssertionService from '../Assertion.service';
import AssertionResultMock from '../../models/__mocks__/AssertionResult.mock';
import AssertionSpanResultMock from '../../models/__mocks__/AssertionSpanResult.mock';

describe('AssertionService', () => {
  describe('getSpanIds', () => {
    it('should return a list of unique span ids', () => {
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
      const result = AssertionService.getSpanIds([assertionResult]);

      expect(result).toEqual([id1, id2, id3]);
    });
  });
});
