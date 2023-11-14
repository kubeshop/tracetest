import {TriggerTypes} from 'constants/Test.constants';
import {TDraftTestForm} from 'types/Test.types';
import Rest from 'components/TestPlugins/Forms/Rest';
import Grpc from 'components/TestPlugins/Forms/Grpc';
import TraceID from 'components/TestPlugins/Forms/TraceID';
import Kafka from 'components/TestPlugins/Forms/Kafka';

const FormFactoryMap = {
  [TriggerTypes.http]: Rest,
  [TriggerTypes.grpc]: Grpc,
  [TriggerTypes.traceid]: TraceID,
  [TriggerTypes.kafka]: Kafka,
};

export interface IFormProps {
  form: TDraftTestForm;
}

interface IProps {
  type: TriggerTypes;
}

const FormFactory = ({type}: IProps) => {
  const Component = FormFactoryMap[type];

  return <Component />;
};

export default FormFactory;
