import {TestState as TestStateEnum} from '../../constants/TestRunResult.constants';
import {ITest} from '../../types/Test.types';
import TestState from '../TestState';
import * as S from './TestHeader.styled';

interface ITestHeaderProps {
  test: ITest;
  onBack(): void;
  testState?: TestStateEnum;
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
          <TestState testState={testState} />
        </S.StateContainer>
      )}
    </S.TestHeader>
  );
};

export default TestHeader;
