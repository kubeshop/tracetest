import {TRawTest, TTest} from '../types/Test.types';
import TestDefinition from './TestDefinition.model';

const Test = ({
  id = '',
  name = '',
  description = '',
  definition,
  serviceUnderTest,
  referenceTestRun,
}: TRawTest): TTest => {
  return {
    id,
    name,
    description,
    definition: definition ? TestDefinition(definition) : undefined,
    serviceUnderTest,
    referenceTestRun,
  };
};

export default Test;
