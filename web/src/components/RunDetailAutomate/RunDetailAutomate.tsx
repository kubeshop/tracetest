import Test from 'models/Test.model';
import * as S from './RunDetailAutomate.styled';
import TestDefinition from '../TestDefinition/TestDefinition';

interface IProps {
  test: Test;
}

const RunDetailAutomate = ({test}: IProps) => {
  return (
    <S.Container>
      <S.SectionLeft>
        <TestDefinition test={test} />
      </S.SectionLeft>
      <S.SectionRight>Test run Techniques</S.SectionRight>
    </S.Container>
  );
};

export default RunDetailAutomate;
