import {Button, Typography} from 'antd';
import {FC, useCallback} from 'react';
import {useRunTestMutation} from 'redux/apis/Test.api';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {Steps} from 'components/GuidedTour/testDetailsStepList';
import useGuidedTour from 'hooks/useGuidedTour';
import * as S from './Test.styled';
import TestDetailsTable from './TestDetailsTable';
import TestAnalyticsService from '../../services/Analytics/TestAnalytics.service';
import {ITestRunResult} from '../../types/TestRunResult.types';

const {onRunTest} = TestAnalyticsService;

type TTestDetailsProps = {
  testId: string;
  url?: string;
  onSelectResult: (result: ITestRunResult) => void;
  testResultList: ITestRunResult[];
  isLoading: boolean;
};

const TestDetails: FC<TTestDetailsProps> = ({testId, testResultList, isLoading, onSelectResult, url}) => {
  const [runTest, result] = useRunTestMutation();
  useGuidedTour(GuidedTours.TestDetails);

  const handleRunTest = useCallback(async () => {
    if (testId) {
      onRunTest(testId);
      const testResult = await runTest(testId).unwrap();
      onSelectResult({resultId: testResult.resultId} as ITestRunResult);
    }
  }, [onSelectResult, runTest, testId]);

  return (
    <div>
      <S.TestDetailsHeader>
        <Typography.Title level={5}>{url}</Typography.Title>
        <Button
          onClick={handleRunTest}
          loading={result.isLoading}
          type="primary"
          data-cy="test-details-run-test-button"
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
