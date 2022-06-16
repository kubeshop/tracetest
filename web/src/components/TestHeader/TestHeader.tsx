import {Button} from 'antd';
import {useState} from 'react';

import TestState from 'components/TestState';
import VersionMismatchModal from 'components/VersionMismatchModal/VersionMismatchModal';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TTest} from 'types/Test.types';
import {TTestRunState} from 'types/TestRun.types';
import Info from './Info';
import * as S from './TestHeader.styled';
import RunActionsDropdown from '../RunActionsMenu';

interface IProps {
  executionTime?: number;
  extraContent?: React.ReactElement;
  onBack(): void;
  showInfo: boolean;
  test: TTest;
  testState?: TTestRunState;
  testVersion: number;
  totalSpans?: number;
}

const TestHeader = ({
  executionTime,
  extraContent,
  onBack,
  showInfo,
  test: {name, referenceTestRun, serviceUnderTest, version = 1, id},
  testState,
  testVersion,
  totalSpans,
}: IProps) => {
  const {runTest} = useTestDefinition();
  const {run} = useTestRun();
  const [isConfirmationModalOpen, setIsConfirmationModalOpen] = useState(false);

  const handleRunTestOnClick = () => {
    TestAnalyticsService.onRunTest();
    if (testVersion !== version) {
      setIsConfirmationModalOpen(true);
      return;
    }
    runTest();
  };

  return (
    <S.TestHeader>
      <S.Content>
        <S.BackIcon data-cy="test-header-back-button" onClick={onBack} />
        <div>
          <S.Row>
            <S.TestName data-cy="test-details-name">
              {name} (v{testVersion})
            </S.TestName>
            {showInfo && (
              <Info
                date={referenceTestRun?.createdAt ?? ''}
                executionTime={executionTime ?? 0}
                totalSpans={totalSpans ?? 0}
                traceId={referenceTestRun?.traceId ?? ''}
              />
            )}
          </S.Row>
          <S.TestUrl>
            {serviceUnderTest?.request?.method?.toUpperCase()} - {serviceUnderTest?.request?.url}
          </S.TestUrl>
        </div>
      </S.Content>
      <S.RightSection>
        {extraContent}
        {!extraContent && testState && testState !== TestStateEnum.FINISHED && (
          <S.StateContainer data-cy="test-run-result-status">
            <S.StateText>Test status:</S.StateText>
            <TestState testState={testState} />
          </S.StateContainer>
        )}
        {!extraContent && testState && testState === TestStateEnum.FINISHED && (
          <Button data-cy="run-test-button" ghost onClick={handleRunTestOnClick} type="primary">
            Run Test
          </Button>
        )}
        {run.id && <RunActionsDropdown resultId={run.id} testId={id} isRunView />}
      </S.RightSection>
      <VersionMismatchModal
        description="Running the test will use the latest version of the test."
        currentVersion={testVersion}
        isOpen={isConfirmationModalOpen}
        latestVersion={version}
        okText="Run Test"
        onCancel={() => setIsConfirmationModalOpen(false)}
        onConfirm={() => {
          setIsConfirmationModalOpen(false);
          runTest();
        }}
      />
    </S.TestHeader>
  );
};

export default TestHeader;
