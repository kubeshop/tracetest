import Test from 'models/Test.model';
import * as S from './RunDetailAutomate.styled';
import TestDefinition from '../RunDetailAutomateDefinition';
import RunDetailAutomateMethods from '../RunDetailAutomateMethods/RunDetailAutomateMethods';

interface IProps {
  test: Test;
}

const RunDetailAutomate = ({test}: IProps) => {
  return (
    <S.Container>
      <S.SectionLeft>
        <TestDefinition test={test} />
      </S.SectionLeft>
      <S.SectionRight>
        <RunDetailAutomateMethods />
      </S.SectionRight>
    </S.Container>
  );
};

export default RunDetailAutomate;
