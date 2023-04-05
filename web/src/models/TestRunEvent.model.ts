import {LogLevel, OutputInfoLogLevel, PollingInfoType, TestRunStage} from 'constants/TestRunEvents.constants';
import {Model, TTestEventsSchemas} from 'types/Common.types';
import ConnectionResult from './ConnectionResult.model';

export type TRawTestRunEvent = TTestEventsSchemas['TestRunEvent'];
export type TRawPollingInfo = TTestEventsSchemas['PollingInfo'];
export type TRawOutputInfo = TTestEventsSchemas['OutputInfo'];

type PollingInfo = Model<TRawPollingInfo, {}>;
type OutputInfo = Model<TRawOutputInfo, {}>;
type TestRunEvent = Model<
  TRawTestRunEvent,
  {logLevel: LogLevel; dataStoreConnection?: ConnectionResult; polling?: PollingInfo; outputs?: OutputInfo[]}
>;

function PollingInfo({
  type = PollingInfoType.Periodic,
  isComplete = false,
  periodic = {},
}: TRawPollingInfo): PollingInfo {
  const types = <any>Object.values(PollingInfoType);

  return {
    type: types.includes(type) ? type : PollingInfoType.Periodic,
    isComplete,
    periodic: {
      numberSpans: periodic?.numberSpans ?? 0,
      numberIterations: periodic?.numberIterations ?? 0,
    },
  };
}

function OutputInfo({
  logLevel = OutputInfoLogLevel.Warning,
  message = '',
  outputName = '',
}: TRawOutputInfo): OutputInfo {
  const logLevels = <any>Object.values(OutputInfoLogLevel);

  return {
    logLevel: logLevels.includes(logLevel) ? logLevel : OutputInfoLogLevel.Warning,
    message,
    outputName,
  };
}

function TestRunEvent({
  type = '',
  stage = TestRunStage.Trigger,
  title = '',
  description = '',
  createdAt = '',
  testId = '',
  runId = 0,
  dataStoreConnection = {},
  polling = {},
  outputs = [],
}: TRawTestRunEvent): TestRunEvent {
  const stages = <any>Object.values(TestRunStage);
  const logLevels = Object.values(LogLevel);
  const logLevel = logLevels.find(level => type.toLowerCase().includes(level.toLowerCase()));

  return {
    type,
    stage: stages.includes(stage) ? stage : TestRunStage.Trigger,
    title,
    description,
    createdAt,
    testId,
    runId,
    logLevel: logLevel ?? LogLevel.Info,
    dataStoreConnection: ConnectionResult(dataStoreConnection),
    polling: PollingInfo(polling),
    outputs: outputs.map(output => OutputInfo(output)),
  };
}

export default TestRunEvent;
