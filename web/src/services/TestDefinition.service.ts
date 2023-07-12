import {TRawTestSpecEntry, TTestSpecEntry} from 'models/TestSpecs.model';

const TestDefinitionService = () => ({
  toRaw({selector, assertions, name}: TTestSpecEntry): TRawTestSpecEntry {
    return {
      selector,
      selectorParsed: {query: selector},
      assertions,
      name,
    };
  },
});

export default TestDefinitionService();
