import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {Button, Typography} from 'antd';
import React, {useCallback, useState} from 'react';

import ResultCardList from 'components/RunCardList';
import {useLazyGetRunListQuery} from 'redux/apis/TraceTest.api';
import {TTest} from 'types/Test.types';
import * as S from './TestCard.styled';
import TestCardActions from './TestCardActions';
import TestAnalyticsService from '../../services/Analytics/TestAnalytics.service';
import {useEditTestModal} from '../EditTestModal/EditTestModal.provider';

interface IProps {
  onClick(testId: string): void;
  onDelete(test: TTest): void;
  onRunTest(testId: string): void;
  test: TTest;
}

const TestCard = ({onClick, onDelete, onRunTest, test: {name, serviceUnderTest, id: testId}, test}: IProps) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [loadResultList, {data: resultList = []}] = useLazyGetRunListQuery();
  const {open} = useEditTestModal();

  const onCollapse = useCallback(async () => {
    TestAnalyticsService.onTestCardCollapse();

    if (resultList.length > 0) {
      setIsCollapsed(true);
      return;
    }
    await loadResultList({testId, take: 5});
    setIsCollapsed(true);
  }, [loadResultList, resultList.length, testId]);

  return (
    <S.TestCard $isCollapsed={isCollapsed}>
      <S.InfoContainer
        onClick={async () => {
          if (isCollapsed) {
            setIsCollapsed(false);
            return;
          }
          onCollapse();
        }}
      >
        {isCollapsed ? <DownOutlined /> : <RightOutlined data-cy={`collapse-test-${testId}`} />}
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
              TestAnalyticsService.onRunTest();
              onRunTest(testId);
            }}
          >
            Run Test
          </Button>
        </S.ButtonContainer>
        <TestCardActions testId={testId} onEdit={() => open(test)} onDelete={() => onDelete(test)} />
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

      {isCollapsed && !resultList.length && (
        <S.EmptyStateContainer>
          <S.EmptyStateIcon />
          <Typography.Text disabled>No Runs</Typography.Text>
        </S.EmptyStateContainer>
      )}
    </S.TestCard>
  );
};

export default TestCard;
