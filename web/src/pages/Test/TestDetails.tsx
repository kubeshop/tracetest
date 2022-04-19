import {skipToken} from '@reduxjs/toolkit/dist/query';
import {Button, Typography} from 'antd';
import {FC, useCallback} from 'react';
import {useGetTestResultsQuery, useRunTestMutation} from 'services/TestService';
import {ITestResult, Test} from 'types';
import * as S from './Test.styled';
import TestDetailsTable from './TestDetailsTable';

type TTestDetailsProps = {
  test?: Test;
  onSelectResult: (result: ITestResult) => void;
};

const TestDetails: FC<TTestDetailsProps> = ({onSelectResult, test}) => {
  const {testId, assertions = [], serviceUnderTest} = test || {};

  const {data: testResults = [], isLoading} = useGetTestResultsQuery(testId ?? skipToken, {
    pollingInterval: 5000,
  });
  const [runTest, result] = useRunTestMutation();

  const handleRunTest = useCallback(() => {
    if (testId) runTest(testId);
  }, [runTest, testId]);

  return (
    <>
      <S.TestDetailsHeader>
        <Typography.Title level={5}>{serviceUnderTest?.request.url}</Typography.Title>
        <Button onClick={handleRunTest} loading={result.isLoading} type="primary" ghost>
          Run Test
        </Button>
      </S.TestDetailsHeader>
      <TestDetailsTable
        assertionList={assertions}
        isLoading={isLoading}
        onSelectResult={onSelectResult}
        testResultList={testResults}
      />
    </>
  );
};

export default TestDetails;
