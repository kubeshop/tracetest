import {Button} from 'antd';
import {useState} from 'react';

import RunActionsMenu from 'components/RunActionsMenu';
import TestState from 'components/TestState';
import TraceActions from 'components/TraceActions';
import VersionMismatchModal from 'components/VersionMismatchModal';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import * as S from './RunDetailLayout.styled';

interface IProps {
  testId: string;
  testVersion: number;
}

const HeaderRight = ({testId, testVersion}: IProps) => {
  const {isDraftMode, runTest} = useTestDefinition();
  const {run} = useTestRun();
  const [isConfirmationModalOpen, setIsConfirmationModalOpen] = useState(false);
  const state = run.state;

  const handleRunTestOnClick = () => {
    TestAnalyticsService.onRunTest();
    if (run.testVersion !== testVersion) {
      setIsConfirmationModalOpen(true);
      return;
    }
    runTest();
  };

  return (
    <>
      <S.Section $justifyContent="flex-end">
        {isDraftMode && <TraceActions />}
        {!isDraftMode && state && state !== TestStateEnum.FINISHED && (
          <S.StateContainer data-cy="test-run-result-status">
            <S.StateText>Test status:</S.StateText>
            <TestState testState={state} />
          </S.StateContainer>
        )}
        {!isDraftMode && state && state === TestStateEnum.FINISHED && (
          <Button data-cy="run-test-button" ghost onClick={handleRunTestOnClick} type="primary">
            Run Test
          </Button>
        )}
        <RunActionsMenu isRunView resultId={run.id} testId={testId} testVersion={testVersion} />
      </S.Section>

      <VersionMismatchModal
        description="Running the test will use the latest version of the test."
        currentVersion={run.testVersion}
        isOpen={isConfirmationModalOpen}
        latestVersion={testVersion}
        okText="Run Test"
        onCancel={() => setIsConfirmationModalOpen(false)}
        onConfirm={() => {
          setIsConfirmationModalOpen(false);
          runTest();
        }}
      />
    </>
  );
};

export default HeaderRight;
