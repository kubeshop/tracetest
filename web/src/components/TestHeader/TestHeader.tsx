import {TestState} from '../../constants/TestRunResult.constants';
import {ITest} from '../../types/Test.types';
import TestStateBadge from '../TestStateBadge';
import * as S from './TestHeader.styled';

interface ITestHeaderProps {
  test: ITest;
  onBack(): void;
  testState?: TestState;
}

const TestHeader: React.FC<ITestHeaderProps> = ({test: {name, serviceUnderTest}, onBack, testState}) => {
  return (
    <S.TestHeader>
      <S.Content>
        <S.BackIcon data-cy="test-header-back-button" onClick={onBack} />
        <S.TestName data-cy="test-details-name">{name}</S.TestName>
        <S.TestUrl>
          {serviceUnderTest.request.method.toUpperCase()} - {serviceUnderTest.request.url}
        </S.TestUrl>
      </S.Content>
      {testState && (
        <S.StateContainer data-cy="test-run-result-status">
          <S.StateText>Test status:</S.StateText>
          <TestStateBadge style={{fontSize: 16}} testState={testState} />
        </S.StateContainer>
      )}
    </S.TestHeader>
  );
};

export default TestHeader;
