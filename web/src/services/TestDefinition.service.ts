import {TTestDefinitionEntry} from '../types/TestDefinition.types';

const TestDefinitionService = () => ({
  toRaw({selector, assertionList}: TTestDefinitionEntry) {

    return {
      selector,
      assertions: assertionList,
    };
  },
});

export default TestDefinitionService();
