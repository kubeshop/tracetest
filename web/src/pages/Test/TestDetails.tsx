import {Button, Typography} from 'antd';
import {FC, useCallback} from 'react';
import {useRunTestMutation} from 'gateways/Test.gateway';
import GuidedTourService, {GuidedTours} from 'entities/GuidedTour/GuidedTour.service';
import {Steps} from 'components/GuidedTour/testDetailsStepList';
import useGuidedTour from 'hooks/useGuidedTour';
import * as S from './Test.styled';
import TestDetailsTable from './TestDetailsTable';
import TestAnalyticsService from '../../entities/Analytics/TestAnalytics.service';
import { TTestRunResult } from '../../entities/TestRunResult/TestRunResult.types';

const {onRunTest} = TestAnalyticsService;

type TTestDetailsProps = {
  testId: string;
  url?: string;
  onSelectResult: (result: TTestRunResult) => void;
  testResultList: TTestRunResult[];
  isLoading: boolean;
};

const TestDetails: FC<TTestDetailsProps> = ({testId, testResultList, isLoading, onSelectResult, url}) => {
  const [runTest, result] = useRunTestMutation();
  useGuidedTour(GuidedTours.TestDetails);

  const handleRunTest = useCallback(async () => {
    if (testId) {
      onRunTest(testId);
      const testResult = await runTest(testId).unwrap();
      onSelectResult({resultId: testResult.resultId} as TTestRunResult);
    }
  }, [onSelectResult, runTest, testId]);

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
