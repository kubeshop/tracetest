import {TRawTestDefinition, TTestDefinition} from '../types/TestDefinition.types';
import Assertion from './Assertion.model';

const TestDefinition = ({specs = []}: TRawTestDefinition): TTestDefinition => {
  return {
    definitionList: specs.map(({selector: {query = ''} = {}, assertions = []}) => ({
      isDraft: false,
      isDeleted: false,
      selector: query,
      assertionList: assertions.map(rawAssertion => Assertion(rawAssertion)),
    })),
  };
};

export default TestDefinition;
