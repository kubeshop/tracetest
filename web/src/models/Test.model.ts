import {TRawTest, TTest} from '../types/Test.types';
import TestDefinition from './TestDefinition.model';
import Trigger from './Trigger.model';

const Test = ({
  id = '',
  name = '',
  description = '',
  definition,
  version = 1,
  serviceUnderTest: rawTrigger,
}: TRawTest): TTest => {
  return {
    id,
    name,
    version,
    description,
    definition: TestDefinition(definition || {}),
    trigger: Trigger(rawTrigger || {}),
  };
};

export default Test;
