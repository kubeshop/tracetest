import {TRawTest, TTest} from 'types/Test.types';
import TestSpecs from './TestSpecs.model';
import Trigger from './Trigger.model';

const Test = ({
  id = '',
  name = '',
  description = '',
  specs,
  version = 1,
  serviceUnderTest: rawTrigger,
}: TRawTest): TTest => ({
  id,
  name,
  version,
  description,
  definition: TestSpecs(specs || {}),
  trigger: Trigger(rawTrigger || {}),
});

export default Test;
