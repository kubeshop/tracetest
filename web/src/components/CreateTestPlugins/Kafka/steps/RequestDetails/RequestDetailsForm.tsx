import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import KeyValueListInput from 'components/KeyValueListInput';
import RequestDetailsBrokerUrlInput from './RequestDetailsBrokerUrlInput';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsTopicInput from './RequestDetailsTopicInput';
import RequestDetailsMessageKeyInput from './RequestDetailsMessageKeyInput';
import RequestDetailsMessageValueInput from './RequestDetailsMessageValueInput';
import SSLVerification from './SSLVerification';

const RequestDetailsForm = () => {
  return (
    <S.InputContainer>
      <RequestDetailsBrokerUrlInput />
      <RequestDetailsAuthInput />
      <RequestDetailsTopicInput />
      <KeyValueListInput name='headers' label='Message Headers' addButtonLabel='Add Header' keyPlaceholder='Header Key' valuePlaceholder='Header Value' />
      <RequestDetailsMessageKeyInput />
      <RequestDetailsMessageValueInput />
      <SSLVerification />
    </S.InputContainer>
  );
};

export default RequestDetailsForm;
