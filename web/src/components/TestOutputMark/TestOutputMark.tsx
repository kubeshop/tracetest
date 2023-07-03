import TestOutput from 'models/TestOutput.model';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import * as S from './TestOutputMark.styled';

interface IProps {
  outputs: TestOutput[];
}

const TestOutputMark = ({outputs}: IProps) => {
  const {onSelectedOutputs} = useTestOutput();

  return (
    <S.Container onClick={() => onSelectedOutputs(outputs)}>
      <S.Text>O</S.Text>
    </S.Container>
  );
};

export default TestOutputMark;
