import PollingForm from './PollingForm';
import * as S from '../common/Settings.styled';

const Polling = () => (
  <S.Container>
    <S.Description>
      Tracetest uses polling to gather the distributed trace associated with each test run. You can adjust the polling
      values below:
    </S.Description>
    <S.FormContainer>
      <PollingForm />
    </S.FormContainer>
  </S.Container>
);

export default Polling;
