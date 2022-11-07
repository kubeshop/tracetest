import {DownOutlined, RightOutlined} from '@ant-design/icons';
import {useMemo} from 'react';

import {useLazyGetRunListQuery} from 'redux/apis/TraceTest.api';
import {ResourceType} from 'types/Resource.type';
import {TTest} from 'types/Test.types';
import {TTestRun} from 'types/TestRun.types';
import ResultCardList from '../RunCardList/RunCardList';
import * as S from './ResourceCard.styled';
import ResourceCardActions from './ResourceCardActions';
import ResourceCardRuns from './ResourceCardRuns';
import ResourceCardSummary from './ResourceCardSummary';
import useRuns from './useRuns';

interface IProps {
  onDelete(id: string, name: string, type: ResourceType): void;
  onRun(id: string, type: ResourceType): void;
  onViewAll(id: string, type: ResourceType): void;
  test: TTest;
}

const TestCard = ({onDelete, onRun, onViewAll, test}: IProps) => {
  const queryParams = useMemo(() => ({take: 5, testId: test.id}), [test.id]);
  const {isCollapsed, isLoading, list, onClick} = useRuns<TTestRun, {testId: string}>(
    useLazyGetRunListQuery,
    queryParams
  );

  return (
    <S.Container $type={ResourceType.Test}>
      <S.TestContainer onClick={onClick}>
        {isCollapsed ? <RightOutlined data-cy={`collapse-test-${test.id}`} /> : <DownOutlined />}
        <S.Box $type={ResourceType.Test}>
          <S.BoxTitle level={2}>{test.summary.runs}</S.BoxTitle>
        </S.Box>
        <S.TitleContainer>
          <S.Title level={3}>{test.name}</S.Title>
          <S.Text>
            {test.trigger.method} • {test.trigger.entryPoint}
          </S.Text>
        </S.TitleContainer>

        <ResourceCardSummary summary={test.summary} />

        <S.Row>
          <S.RunButton
            type="primary"
            ghost
            data-cy={`test-run-button-${test.id}`}
            onClick={event => {
              event.stopPropagation();
              onRun(test.id, ResourceType.Test);
            }}
          >
            Run
          </S.RunButton>
          <ResourceCardActions id={test.id} onDelete={() => onDelete(test.id, test.name, ResourceType.Test)} />
        </S.Row>
      </S.TestContainer>

      <ResourceCardRuns
        hasMoreRuns={list.length === 5}
        hasRuns={Boolean(list.length)}
        isCollapsed={isCollapsed}
        isLoading={isLoading}
        onViewAll={() => onViewAll(test.id, ResourceType.Test)}
      >
        <ResultCardList testId={test.id} resultList={list} />
      </ResourceCardRuns>
    </S.Container>
  );
};

export default TestCard;
