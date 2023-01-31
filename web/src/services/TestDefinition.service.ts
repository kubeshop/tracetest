import {TRawTestSpecEntry, TTestSpecEntry} from 'models/TestSpecs.model';

const TestDefinitionService = () => ({
  toRaw({selector, assertions}: TTestSpecEntry): TRawTestSpecEntry {
    return {
      selector: {query: selector},
      assertions,
    };
  },
});

export default TestDefinitionService();
