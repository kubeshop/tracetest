import {TRawTest, TTest} from '../types/Test.types';
import TestDefinition from './TestDefinition.model';

const Test = ({id = '', name = '', description = '', definition, version = 1, serviceUnderTest}: TRawTest): TTest => {
  return {
    id,
    name,
    version,
    description,
    definition: TestDefinition(definition || {}),
    serviceUnderTest,
  };
};

export default Test;
