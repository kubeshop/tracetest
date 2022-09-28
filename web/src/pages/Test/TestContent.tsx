import {useNavigate} from 'react-router-dom';

import TestHeader from 'components/TestHeader';
import {useTest} from 'providers/Test/Test.provider';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import * as S from './Test.styled';
import TestDetails from './TestDetails';

const TestContent = () => {
  const navigate = useNavigate();
  const {test} = useTest();

  return (
    <S.Wrapper $isWhite={!ExperimentalFeature.isEnabled('transactions')}>
      <TestHeader onBack={() => navigate('/')} test={test} />
      <TestDetails testId={test.id} />
    </S.Wrapper>
  );
};

export default TestContent;
