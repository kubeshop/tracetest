import {TTestOutput} from 'types/TestOutput.types';
import {useTestOutput} from '../../providers/TestOutput/TestOutput.provider';
import * as S from './TestOutputMark.styled';

interface IProps {
  className?: string;
  outputs: TTestOutput[];
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
