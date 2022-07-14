import {useCallback} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import {useGetTestByIdQuery} from 'redux/apis/TraceTest.api';
import {TTestRun} from 'types/TestRun.types';
import TestHeader from 'components/TestHeader';
import TestCardActions from 'components/TestCard/TestCardActions';
import * as S from './Test.styled';
import TestDetails from './TestDetails';
import {useMenuDeleteCallback} from '../Home/useMenuDeleteCallback';

const TestContent = () => {
  const navigate = useNavigate();
  const {testId = ''} = useParams();
  const {data: test} = useGetTestByIdQuery({testId});

  const handleSelectTestResult = useCallback(
    (result: TTestRun) => {
      navigate(`/test/${testId}/run/${result.id}`);
    },
    [navigate, testId]
  );
  const onDelete = useMenuDeleteCallback();

  // TODO: Add proper loading states
  return test ? (
    <>
      <TestHeader
        onBack={() => navigate('/')}
        showInfo={false}
        test={test}
        testVersion={test.version}
        extraContent={<TestCardActions testId={testId} onDelete={() => onDelete(test)} />}
      />
      <S.Wrapper>
        <TestDetails testId={test.id} onSelectResult={handleSelectTestResult} />
      </S.Wrapper>
    </>
  ) : null;
};

export default TestContent;
