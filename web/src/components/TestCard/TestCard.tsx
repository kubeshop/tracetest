import {RightOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import {ITest} from '../../types/Test.types';
import * as S from './TestCard.styled';
import TestCardActions from './TestCardActions';

interface ITestCardProps {
  test: ITest;
  onClick(testId: string): void;
  onDelete(test: ITest): void;
  onRunTest(testId: string): void;
}

const TestCard: React.FC<ITestCardProps> = ({
  test: {name, serviceUnderTest, testId},
  test,
  onClick,
  onDelete,
  onRunTest,
}) => {
  return (
    <S.TestCard onClick={() => onClick(testId)} data-cy="test-card">
      <RightOutlined />
      <S.TextContainer>
        <S.NameText>{name}</S.NameText>
      </S.TextContainer>
      <S.TextContainer>
        <S.Text>{serviceUnderTest.request.method}</S.Text>
      </S.TextContainer>
      <S.TextContainer data-cy={`test-url-${testId}`}>
        <S.Text>{serviceUnderTest.request.url}</S.Text>
      </S.TextContainer>
      <S.TextContainer />
      <S.ButtonContainer>
        <Button
          type="primary"
          ghost
          data-cy="test-run-button"
          onClick={event => {
            event.stopPropagation();
            onRunTest(testId);
          }}
        >
          Run Test
        </Button>
      </S.ButtonContainer>
      <TestCardActions testId={testId} onDelete={() => onDelete(test)} />
    </S.TestCard>
  );
};

export default TestCard;
