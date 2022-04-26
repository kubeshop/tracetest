import {Button, Typography} from 'antd';
import {FC, useCallback} from 'react';
import {useRunTestMutation} from 'redux/services/TestService';
import {TestId, TestRunResult} from 'types';
import GuidedTourService, {GuidedTours} from 'services/GuidedTourService';
import {Steps} from 'components/GuidedTour/testDetailsStepList';
import useGuidedTour from 'hooks/useGuidedTour';
import * as S from './Test.styled';
import TestDetailsTable from './TestDetailsTable';
import useTestAnalytics from './useTest.analytics';

type TTestDetailsProps = {
  testId: TestId;
  url?: string;
  onSelectResult: (result: TestRunResult) => void;
  testResultList: TestRunResult[];
  isLoading: boolean;
};

const TestDetails: FC<TTestDetailsProps> = ({testId, testResultList, isLoading, onSelectResult, url}) => {
  const [runTest, result] = useRunTestMutation();
  useGuidedTour(GuidedTours.TestDetails);
  const {onRunTest} = useTestAnalytics();

  const handleRunTest = useCallback(async () => {
    if (testId) {
      onRunTest(testId);
      const testResult = await runTest(testId).unwrap();
      onSelectResult({resultId: testResult.resultId} as TestRunResult);
    }
  }, [onRunTest, onSelectResult, runTest, testId]);

  return (
    <div style={{height: 'calc(100vh - 250px)'}}>
      <S.TestDetailsHeader>
        <Typography.Title level={5}>{url}</Typography.Title>
        <Button
          onClick={handleRunTest}
          loading={result.isLoading}
          type="primary"
          ghost
          data-tour={GuidedTourService.getStep(GuidedTours.TestDetails, Steps.RunTest)}
        >
          Run Test
        </Button>
      </S.TestDetailsHeader>
      <TestDetailsTable isLoading={isLoading} onSelectResult={onSelectResult} testResultList={testResultList} />
    </div>
  );
};

export default TestDetails;
