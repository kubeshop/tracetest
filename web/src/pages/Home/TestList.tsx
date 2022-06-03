import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import TestCard from '../../components/TestCard';
import {useGetTestListQuery, useRunTestMutation} from '../../redux/apis/TraceTest.api';
import * as S from './Home.styled';
import {useMenuDeleteCallback} from './useMenuDeleteCallback';

const {onTestClick} = HomeAnalyticsService;

const TestList = () => {
  const navigate = useNavigate();
  const {data: testList = []} = useGetTestListQuery();
  const [runTest] = useRunTestMutation();

  const onClick = useCallback(
    (testId: string) => {
      onTestClick(testId);
      navigate(`/test/${testId}`);
    },
    [navigate]
  );

  const onRunTest = useCallback(
    async (testId: string) => {
      if (testId) {
        TestAnalyticsService.onRunTest();
        const testRun = await runTest({testId}).unwrap();
        navigate(`/test/${testId}/run/${testRun.id}`);
      }
    },
    [navigate, runTest]
  );

  const onDelete = useMenuDeleteCallback();

  return (
    <S.TestListContainer data-cy="test-list">
      {testList?.map(test => (
        <TestCard test={test} onClick={onClick} onDelete={onDelete} onRunTest={onRunTest} key={test.id} />
      ))}
    </S.TestListContainer>
  );
};

export default TestList;
