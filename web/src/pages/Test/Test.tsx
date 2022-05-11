import {useCallback} from 'react';
import {skipToken} from '@reduxjs/toolkit/dist/query';
import {withTracker} from 'ga-4-react';
import {useNavigate, useParams} from 'react-router-dom';
import {useGetTestByIdQuery, useGetResultListQuery} from 'redux/apis/Test.api';
import Layout from 'components/Layout';

import * as S from './Test.styled';
import TestDetails from './TestDetails';
import {ITestRunResult} from '../../types/TestRunResult.types';
import TestHeader from '../../components/TestHeader';

const TestPage: React.FC = () => {
  const navigate = useNavigate();
  const {id} = useParams();
  const {data: test} = useGetTestByIdQuery(id as string);
  const {data: testResultList = [], isLoading} = useGetResultListQuery(id ?? skipToken, {
    pollingInterval: 5000,
  });

  const handleSelectTestResult = useCallback(
    (result: ITestRunResult) => {
      navigate(`/test/${id}/result/${result.resultId}`);
    },
    [id, navigate]
  );

  // TODO: Add proper loading states
  return test ? (
    <Layout>
      <TestHeader test={test} onBack={() => navigate('/')} />
      <S.Wrapper>
        <TestDetails
          testResultList={testResultList}
          isLoading={isLoading}
          testId={id!}
          onSelectResult={handleSelectTestResult}
        />
      </S.Wrapper>
    </Layout>
  ) : null;
};

export default withTracker(TestPage);
