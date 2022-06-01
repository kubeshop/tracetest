import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import React, {useCallback, useState} from 'react';
import {useLazyGetRunListQuery} from '../../redux/apis/TraceTest.api';
import {TTest} from '../../types/Test.types';
import ResultCardList from '../RunCardList';
import * as S from './TestCard.styled';
import TestCardActions from './TestCardActions';

interface TTestCardProps {
  test: TTest;

  onClick(testId: string): void;

  onDelete(test: TTest): void;

  onRunTest(testId: string): void;
}

const TestCard: React.FC<TTestCardProps> = ({
  test: {name, serviceUnderTest, id: testId},
  test,
  onClick,
  onDelete,
  onRunTest,
}) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [loadResultList, {data: resultList = []}] = useLazyGetRunListQuery();

  const onCollapse = useCallback(async () => {
    if (!resultList.length) {
      const list = await loadResultList({testId, take: 5}).unwrap();

      if (list.length) setIsCollapsed(true);
    } else {
      setIsCollapsed(true);
    }
  }, [loadResultList, resultList.length, testId]);

  return (
    <S.TestCard $isCollapsed={isCollapsed}>
      <S.InfoContainer
        onClick={async () => {
          if (isCollapsed) {
            setIsCollapsed(false);
            return;
          }
          await onCollapse();
        }}
      >
        {isCollapsed ? <DownOutlined /> : <RightOutlined data-cy={`collapse-test-${testId}`} onClick={onCollapse} />}
        <S.TextContainer>
          <S.NameText>{name}</S.NameText>
        </S.TextContainer>
        <S.TextContainer>
          <S.Text>{serviceUnderTest?.request?.method}</S.Text>
        </S.TextContainer>
        <S.TextContainer data-cy={`test-url-${testId}`}>
          <S.Text>{serviceUnderTest?.request?.url}</S.Text>
        </S.TextContainer>
        <S.TextContainer />
        <S.ButtonContainer>
          <Button
            type="primary"
            ghost
            data-cy={`test-run-button-${testId}`}
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
          <ResultCardList testId={testId} resultList={resultList} />
          {resultList.length === 5 && (
            <S.TestDetails>
              <S.TestDetailsLink data-cy="test-details-link" onClick={() => onClick(testId)}>
                Explore all test details
              </S.TestDetailsLink>
            </S.TestDetails>
          )}
        </S.ResultListContainer>
      )}
    </S.TestCard>
  );
};

export default TestCard;
