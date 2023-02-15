import AnalyticsForm from './AnalyticsForm';
import * as S from '../common/Settings.styled';

const Analytics = () => (
  <S.Container>
    <S.Description>
      To improve the end-user experience, Tracetest collects anonymous telemetry data about usage. Participation is
      optional, and you can opt-out by turning analytics off. See telemetry documentation for more information.
    </S.Description>
    <S.FormContainer>
      <AnalyticsForm />
    </S.FormContainer>
  </S.Container>
);

export default Analytics;
