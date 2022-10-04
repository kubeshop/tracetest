import {TRawTestSpecEntry, TTestSpecEntry} from 'types/TestSpecs.types';
import AssertionService from './Assertion.service';

const TestDefinitionService = () => ({
  toRaw({selector, assertions}: TTestSpecEntry): TRawTestSpecEntry {
    return {
      selector: {query: selector},
      assertions,
    };
  },
  formatExpectedField(rawTestSpecs: TRawTestSpecEntry[]) {
    return rawTestSpecs.map(spec => ({
      ...spec,
      assertions: spec.assertions.map(assertion => ({
        ...assertion,
        expected: AssertionService.extractExpectedString(assertion.expected),
      })),
    }));
  },
});

export default TestDefinitionService();
