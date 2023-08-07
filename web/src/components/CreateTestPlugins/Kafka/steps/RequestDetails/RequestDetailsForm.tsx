import {IKafkaValues, TDraftTestForm} from 'types/Test.types';

interface IProps {
  form: TDraftTestForm<IKafkaValues>;
}

// TODO daniel
const RequestDetailsForm = ({form}: IProps) => {
  return (
    <p>Form: {form.getFieldValue("sslVerification")}</p>
  );
};

export default RequestDetailsForm;
