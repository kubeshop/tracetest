import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import * as S from './RunDetailAutomate.styled';

interface IProps {
  run: TestRun;
  runEvents: TestRunEvent[];
  testId: string;
}

export enum VisualizationType {
  Dag,
  Timeline,
}

const RunDetailAutomate = ({run, runEvents, testId}: IProps) => {
  console.log('@@@props', run, runEvents, testId);

  return (
    <S.Container>
      <S.SectionLeft>Test Definition</S.SectionLeft>
      <S.SectionRight>Test run Techniques</S.SectionRight>
    </S.Container>
  );
};

export default RunDetailAutomate;
