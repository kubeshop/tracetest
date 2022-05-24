import faker from '@faker-js/faker';
import AssertionService from '../Assertion.service';
import AssertionResultMock from '../../models/__mocks__/AssertionResult.mock';
import AssertionSpanResultMock from '../../models/__mocks__/AssertionSpanResult.mock';

describe('AssertionService', () => {
  describe('getSpanCount', () => {
    test('should returning the number of spans', () => {
      const id = faker.datatype.uuid();
      const assertionResult = AssertionResultMock.model({
        spanResults: [
          AssertionSpanResultMock.raw({
            spanId: id,
          }),
          AssertionSpanResultMock.raw({
            spanId: id,
          }),
          AssertionSpanResultMock.raw(),
          AssertionSpanResultMock.raw(),
        ],
      });

      const result = AssertionService.getSpanCount([assertionResult]);

      expect(result).toBe(3);
    });
  });
});
