import {useSetIsCollapsedCallback} from 'components/ResizableDrawer/useSetIsCollapsedCallback';
import TestState from 'components/TestState';
import {TTest} from 'types/Test.types';
import {TTestRunState} from 'types/TestRun.types';
import Info from './Info';
import * as S from './TestHeader.styled';

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
  test: {name, referenceTestRun, serviceUnderTest},
  test,
  testState,
  testVersion,
  totalSpans,
}: IProps) => {
  const onClick = useSetIsCollapsedCallback();

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
      {extraContent}
      {testState && !extraContent && (
        <S.StateContainer onClick={onClick} data-cy="test-run-result-status">
          <S.StateText>Test status:</S.StateText>
          <TestState testState={testState} />
        </S.StateContainer>
      )}
    </S.TestHeader>
  );
};

export default TestHeader;
