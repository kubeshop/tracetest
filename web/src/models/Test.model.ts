import {TRawTest, TTest} from 'types/Test.types';
import TestOutput from './TestOutput.model';
import TestSpecs from './TestSpecs.model';
import TestSummary from './TestSummary.model';
import Trigger from './Trigger.model';

const Test = ({
  id = '',
  name = '',
  description = '',
  specs,
  version = 1,
  serviceUnderTest: rawTrigger,
  summary = {},
  outputs = [],
}: TRawTest): TTest => ({
  id,
  name,
  version,
  description,
  definition: TestSpecs(specs || {}),
  trigger: Trigger(rawTrigger || {}),
  summary: TestSummary(summary),
  outputs: outputs?.map(rawOutput => TestOutput(rawOutput)),
});

export default Test;
