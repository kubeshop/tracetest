import {TriggerTypes} from '../../constants/Test.constants';
import {TTestRunState} from '../../types/TestRun.types';
import {TTriggerResult} from '../../types/Test.types';
import RunDetailTriggerData from './RunDetailTriggerData';
import RunDetailTriggerResponse from './RunDetailTriggerResponse';

export interface IPropsComponent {
  state: TTestRunState;
  triggerResult?: TTriggerResult;
  triggerTime?: number;
}

const ComponentMap: Record<TriggerTypes, (props: IPropsComponent) => React.ReactElement> = {
  [TriggerTypes.http]: RunDetailTriggerResponse,
  [TriggerTypes.grpc]: RunDetailTriggerResponse,
  [TriggerTypes.traceid]: RunDetailTriggerData,
};

interface IProps extends IPropsComponent {
  type: TriggerTypes;
}

const RunDetailTriggerResponseFactory = ({type, ...props}: IProps) => {
  const Component = ComponentMap[type];

  return <Component {...props} />;
};

export default RunDetailTriggerResponseFactory;
