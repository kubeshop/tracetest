import {TRawTestDefinitionEntry, TTestDefinitionEntry} from '../types/TestDefinition.types';

const TestDefinitionService = () => ({
  toRaw({selector, assertionList}: TTestDefinitionEntry): TRawTestDefinitionEntry {

    return {
      selector,
      assertions: assertionList,
    };
  },
});

export default TestDefinitionService();
