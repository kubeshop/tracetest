import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import {useCallback, useState} from 'react';
import {useLazyGetResultListQuery} from '../../redux/apis/Test.api';
import {ITest} from '../../types/Test.types';
import ResultCardList from '../ResultCardList';
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
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [loadResultList, {data: resultList = []}] = useLazyGetResultListQuery();

  const onCollapse = useCallback(async () => {
    if (!resultList.length) {
      const list = await loadResultList({testId, take: 5}).unwrap();

      if (list.length) setIsCollapsed(true);
    } else {
      setIsCollapsed(true);
    }
  }, [loadResultList, resultList, testId]);

  return (
    <S.TestCard isCollapsed={isCollapsed}>
      <S.InfoContainer>
        {isCollapsed ? (
          <DownOutlined onClick={() => setIsCollapsed(false)} />
        ) : (
          <RightOutlined data-cy="collapse-test" onClick={onCollapse} />
        )}
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
      </S.InfoContainer>

      {isCollapsed && Boolean(resultList.length) && (
        <S.ResultListContainer>
          <ResultCardList resultList={resultList} />
          <S.TestDetails>
            <S.TestDetailsLink data-cy="test-details-link" onClick={() => onClick(testId)}>
              Explore all test details
            </S.TestDetailsLink>
          </S.TestDetails>
        </S.ResultListContainer>
      )}
    </S.TestCard>
  );
};

export default TestCard;
