import {TConnectionResult} from 'types/Config.types';
import TestConnectionStep from './TestConnectionStep';
import * as S from './TestConnectionNotification.styled';

interface IProps {
  result: TConnectionResult;
}

const TestConnectionNotification = ({result: {authentication, connectivity, fetchTraces}}: IProps) => {
  return (
    <S.Container>
      <TestConnectionStep step={connectivity} />
      <TestConnectionStep step={authentication} />
      <TestConnectionStep step={fetchTraces} />
    </S.Container>
  );
};

export default TestConnectionNotification;
