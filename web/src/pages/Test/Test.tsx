import {useCallback} from 'react';
import {withTracker} from 'ga-4-react';
import {useNavigate, useParams} from 'react-router-dom';
import {useGetTestByIdQuery} from 'redux/apis/TraceTest.api';
import Layout from 'components/Layout';

import * as S from './Test.styled';
import TestDetails from './TestDetails';
import {TTestRun} from '../../types/TestRun.types';
import TestHeader from '../../components/TestHeader';

const TestPage: React.FC = () => {
  const navigate = useNavigate();
  const {testId = ''} = useParams();
  const {data: test} = useGetTestByIdQuery({testId});

  const handleSelectTestResult = useCallback(
    (result: TTestRun) => {
      navigate(`/test/${testId}/result/${result.id}`);
    },
    [navigate, testId]
  );

  // TODO: Add proper loading states
  return test ? (
    <Layout>
      <TestHeader test={test} onBack={() => navigate('/')} />
      <S.Wrapper>
        <TestDetails testId={test.id} onSelectResult={handleSelectTestResult} />
      </S.Wrapper>
    </Layout>
  ) : null;
};

export default withTracker(TestPage);
