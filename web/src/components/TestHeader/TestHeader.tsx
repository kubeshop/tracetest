import {TTest} from '../../types/Test.types';
import {TTestRunState} from '../../types/TestRun.types';
import TestState from '../TestState';
import * as S from './TestHeader.styled';

interface TTestHeaderProps {
  test: TTest;
  onBack(): void;
  testState?: TTestRunState;
}

const TestHeader: React.FC<TTestHeaderProps> = ({test: {name, serviceUnderTest}, onBack, testState}) => {
  return (
    <S.TestHeader>
      <S.Content>
        <S.BackIcon data-cy="test-header-back-button" onClick={onBack} />
        <S.TestName data-cy="test-details-name">{name}</S.TestName>
        <S.TestUrl>
          {serviceUnderTest?.request?.method?.toUpperCase()} - {serviceUnderTest?.request?.url}
        </S.TestUrl>
      </S.Content>
      {testState && (
        <S.StateContainer data-cy="test-run-result-status">
          <S.StateText>Test status:</S.StateText>
          <TestState testState={testState} />
        </S.StateContainer>
      )}
    </S.TestHeader>
  );
};

export default TestHeader;
