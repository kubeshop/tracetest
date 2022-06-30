import SelectorService from '../services/Selector.service';
import {TRawTestDefinition, TTestDefinition} from '../types/TestDefinition.types';
import Assertion from './Assertion.model';

const TestDefinition = ({definitions = []}: TRawTestDefinition): TTestDefinition => {
  return {
    definitionList: definitions.map(({selector: {query = ''} = {}, assertions = []}) => ({
      isDraft: false,
      isDeleted: false,
      isAdvancedSelector: SelectorService.getIsAdvancedSelector(query),
      selector: query,
      assertionList: assertions.map(rawAssertion => Assertion(rawAssertion)),
    })),
  };
};

export default TestDefinition;
