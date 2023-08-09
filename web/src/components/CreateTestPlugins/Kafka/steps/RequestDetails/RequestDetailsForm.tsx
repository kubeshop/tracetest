import {IKafkaValues, TDraftTestForm} from 'types/Test.types';
import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import RequestDetailsBrokerUrlInput from './RequestDetailsBrokerUrlInput';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsTopicInput from './RequestDetailsTopicInput';
import RequestDetailsHeadersInput from './RequestDetailsHeadersInput';
import RequestDetailsMessageKeyInput from './RequestDetailsMessageKeyInput';
import RequestDetailsMessageValueInput from './RequestDetailsMessageValueInput';
import SSLVerification from './SSLVerification';

interface IProps {
  form: TDraftTestForm<IKafkaValues>;
}

const RequestDetailsForm = ({form}: IProps) => {
  return (
    <S.InputContainer>
      <RequestDetailsBrokerUrlInput />
      <RequestDetailsAuthInput />
      <RequestDetailsTopicInput />
      <RequestDetailsHeadersInput />
      <RequestDetailsMessageKeyInput />
      <RequestDetailsMessageValueInput />
      <SSLVerification />
    </S.InputContainer>
  );
};

export default RequestDetailsForm;
