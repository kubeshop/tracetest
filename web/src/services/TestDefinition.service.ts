import {TRawTestSpecEntry, TTestSpecEntry} from 'types/TestSpecs.types';

const TestDefinitionService = () => ({
  toRaw({selector, assertions}: TTestSpecEntry): TRawTestSpecEntry {
    return {
      selector: {query: selector},
      assertions,
    };
  },
});

export default TestDefinitionService();
