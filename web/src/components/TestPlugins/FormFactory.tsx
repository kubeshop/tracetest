import Rest from 'components/TestPlugins/Forms/Rest';
import Grpc from 'components/TestPlugins/Forms/Grpc';
import Kafka from 'components/TestPlugins/Forms/Kafka';
import {TriggerTypes} from 'constants/Test.constants';

const FormFactoryMap = {
  [TriggerTypes.http]: Rest,
  [TriggerTypes.grpc]: Grpc,
  [TriggerTypes.kafka]: Kafka,
  [TriggerTypes.traceid]: () => null,
  [TriggerTypes.cypress]: () => null,
  [TriggerTypes.playwright]: () => null,
};

interface IProps {
  type: TriggerTypes;
}

const FormFactory = ({type}: IProps) => {
  const Component = FormFactoryMap[type];

  return <Component />;
};

export default FormFactory;
