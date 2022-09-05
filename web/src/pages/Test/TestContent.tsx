import TestHeader from 'components/TestHeader';
import * as S from './Test.styled';
import TestDetails from './TestDetails';
import {useTest} from '../../providers/Test/Test.provider';

const TestContent = () => {
  const {test} = useTest();

  return (
    <>
      <TestHeader test={test} />
      <S.Wrapper>
        <TestDetails testId={test.id} />
      </S.Wrapper>
    </>
  );
};

export default TestContent;
