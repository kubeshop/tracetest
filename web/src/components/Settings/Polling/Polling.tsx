import DocsBanner from 'components/DocsBanner';
import {TRACE_POLLING_DOCUMENTATION_URL} from 'constants/Common.constants';
import PollingForm from './PollingForm';
import * as S from '../common/Settings.styled';

const Polling = () => (
  <S.Container>
    <S.Title level={2}>Trace Polling</S.Title>
    <S.Description>
      <p>Tracetest uses polling to gather the distributed trace associated with each test run.</p>
      <DocsBanner>
        Need more information about Trace Polling?{' '}
        <a href={TRACE_POLLING_DOCUMENTATION_URL} target="__blank">
          Learn more in our docs
        </a>
      </DocsBanner>
    </S.Description>
    <S.FormContainer>
      <PollingForm />
    </S.FormContainer>
  </S.Container>
);

export default Polling;
