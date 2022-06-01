import {TTest} from '../../types/Test.types';
import {TTestRunState} from '../../types/TestRun.types';
import {useSetIsCollapsedCallback} from '../ResizableDrawer/useSetIsCollapsedCallback';
import TestState from '../TestState';
import * as S from './TestHeader.styled';

interface TTestHeaderProps {
  test: TTest;
  onBack(): void;
  testState?: TTestRunState;
  testVersion: number;
  extraContent?: React.ReactElement;
}

const TestHeader: React.FC<TTestHeaderProps> = ({
  test: {name, serviceUnderTest},
  onBack,
  testState,
  extraContent,
  testVersion,
}) => {
  const onClick = useSetIsCollapsedCallback();
  return (
    <S.TestHeader>
      <S.Content>
        <S.BackIcon data-cy="test-header-back-button" onClick={onBack} />
        <S.TestName data-cy="test-details-name">
          {name} (v{testVersion})
        </S.TestName>
        <S.TestUrl>
          {serviceUnderTest?.request?.method?.toUpperCase()} - {serviceUnderTest?.request?.url}
        </S.TestUrl>
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
