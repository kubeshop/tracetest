import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import * as S from './RunDetailAutomate.styled';
import TestDefinition from '../RunDetailAutomateDefinition';
import RunDetailAutomateMethods from '../RunDetailAutomateMethods/RunDetailAutomateMethods';

interface IProps {
  test: Test;
  run: TestRun;
}

const RunDetailAutomate = ({test, run}: IProps) => (
  <S.Container>
    <S.SectionLeft>
      <TestDefinition test={test} />
    </S.SectionLeft>
    <S.SectionRight>
      <RunDetailAutomateMethods test={test} run={run} />
    </S.SectionRight>
  </S.Container>
);

export default RunDetailAutomate;
