import {Button, Typography} from 'antd';
import {FC, useCallback} from 'react';
import {useRunTestMutation} from 'services/TestService';
import {TestId, TestRunResult} from 'types';
import * as S from './Test.styled';
import TestDetailsTable from './TestDetailsTable';

type TTestDetailsProps = {
  testId: TestId;
  url?: string;
  onSelectResult: (result: TestRunResult) => void;
  testResultList: TestRunResult[];
  isLoading: boolean;
};

const TestDetails: FC<TTestDetailsProps> = ({testId, testResultList, isLoading, onSelectResult, url}) => {
  const [runTest, result] = useRunTestMutation();

  const handleRunTest = useCallback(() => {
    if (testId) runTest(testId);
  }, [runTest, testId]);

  return (
    <>
      <S.TestDetailsHeader>
        <Typography.Title level={5}>{url}</Typography.Title>
        <Button onClick={handleRunTest} loading={result.isLoading} type="primary" ghost>
          Run Test
        </Button>
      </S.TestDetailsHeader>
      <TestDetailsTable isLoading={isLoading} onSelectResult={onSelectResult} testResultList={testResultList} />
    </>
  );
};

export default TestDetails;
