import {TRawTestSpecEntry, TTestSpecEntry} from 'types/TestSpecs.types';

const TestDefinitionService = () => ({
  toRaw({selector, assertions, name}: TTestSpecEntry): TRawTestSpecEntry {
    return {
      selector: {query: selector},
      assertions,
      name,
    };
  },
});

export default TestDefinitionService();
