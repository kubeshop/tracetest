import TestOutput from 'models/TestOutput.model';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import * as S from './TestOutputMark.styled';

interface IProps {
  className?: string;
  outputs: TestOutput[];
}

const TestOutputMark = ({className, outputs}: IProps) => {
  const {onSelectedOutputs} = useTestOutput();

  return (
    <S.Mark onClick={() => onSelectedOutputs(outputs)} className={className}>
      [x]
    </S.Mark>
  );
};

export default TestOutputMark;
