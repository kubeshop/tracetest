import {TConnectionResult} from 'types/Config.types';
import TestConnectionStep from './TestConnectionStep';
import * as S from './TestConnectionNotification.styled';

interface IProps {
  result: TConnectionResult;
}

const TestConnectionNotification = ({result: {endpointLinting, authentication, connectivity, fetchTraces}}: IProps) => {
  return (
    <S.Container>
      <TestConnectionStep step={endpointLinting} title="Port checking" />
      <TestConnectionStep step={connectivity} title="Connectivity" />
      <TestConnectionStep step={authentication} title="Authentication" />
      <TestConnectionStep step={fetchTraces} title="Fetch traces" />
    </S.Container>
  );
};

export default TestConnectionNotification;
