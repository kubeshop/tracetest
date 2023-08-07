import BasicDetailsForm from 'components/CreateTestPlugins/Kafka/steps/RequestDetails/RequestDetailsForm';
import {IKafkaValues, TDraftTestForm} from 'types/Test.types';
import {IFormProps} from '../EditRequestDetails';

const EditRequestDetailsKafka = ({form}: IFormProps) => <BasicDetailsForm form={form as TDraftTestForm<IKafkaValues>} />;

export default EditRequestDetailsKafka;
