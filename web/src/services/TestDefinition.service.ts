import {TRawTestSpecEntry, TTestSpecEntry} from 'types/TestSpecs.types';

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
      assertions: spec.assertions,
    }));
  },
});

export default TestDefinitionService();
