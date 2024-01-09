import {TriggerTypes} from 'constants/Test.constants';
import {load} from 'js-yaml';
import {Model, TTestSchemas} from 'types/Common.types';
import TestOutput from './TestOutput.model';
import TestSpecs from './TestSpecs.model';
import Summary from './Summary.model';
import Trigger from './Trigger.model';

export type TRawTestResource = TTestSchemas['TestResource'];
export type TRawTestResourceList = TTestSchemas['TestResourceList'];
export type TRawTest = TTestSchemas['Test'];
type Test = Model<
  TRawTest,
  {
    definition: TestSpecs;
    serviceUnderTest?: undefined;
    trigger: Trigger;
    specs?: TestSpecs;
    summary: Summary;
    outputs?: TestOutput[];
    createdAt?: string;
  }
>;

const Test = ({spec: rawTest = {}}: TRawTestResource): Test => Test.FromRawTest(rawTest);

Test.FromRawTest = ({
  id = '',
  name = '',
  description = '',
  specs = [],
  version = 1,
  trigger: rawTrigger,
  summary = {},
  outputs = [],
  createdAt = '',
  skipTraceCollection = false,
}: TRawTest): Test => {
  return {
    id,
    name,
    version,
    description,
    createdAt,
    skipTraceCollection,
    definition: TestSpecs({specs}),
    trigger: Trigger(rawTrigger || {}),
    summary: Summary(summary),
    outputs: outputs.map((rawOutput, index) => TestOutput(rawOutput, index)),
  };
};

Test.FromDefinition = (definition: string): Test => {
  const raw: TRawTestResource = load(definition);
  return Test(raw);
};

Test.shouldAllowRun = (triggerType: TriggerTypes): boolean => {
  return ![TriggerTypes.cypress, TriggerTypes.playwright].includes(triggerType);
};

export default Test;
