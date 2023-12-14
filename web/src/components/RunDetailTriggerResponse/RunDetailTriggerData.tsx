import {Tabs} from 'antd';
import {useSearchParams} from 'react-router-dom';
import {StepsID} from 'components/GuidedTour/testRunSteps';
import {TestState} from 'constants/TestRun.constants';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import ResponseVariableSet from './ResponseVariableSet';
import * as S from './RunDetailTriggerResponse.styled';
import {IPropsComponent} from './RunDetailTriggerResponseFactory';
import ResponseMetadata from './ResponseMetadata';

const TABS = {
  VariableSet: 'variable-set',
  Metadata: 'metadata',
} as const;

const RunDetailTriggerData = ({state, triggerTime = 0}: IPropsComponent) => {
  const [query, updateQuery] = useSearchParams();

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>Trigger Data</S.Title>
        <div>
          <S.StatusText>
            Time:{' '}
            <S.StatusSpan $isError={triggerTime > 1000}>
              {state === TestState.CREATED || state === TestState.EXECUTING ? '-' : `${triggerTime}ms`}
            </S.StatusSpan>
          </S.StatusText>
        </div>
      </S.TitleContainer>
      <S.TabsContainer data-tour={StepsID.Response}>
        <Tabs
          defaultActiveKey={query.get('tab') || TABS.VariableSet}
          size="small"
          onChange={newTab => {
            TestRunAnalyticsService.onTriggerResponseTabChange(newTab);
            updateQuery([['tab', newTab]]);
          }}
        >
          <Tabs.TabPane key={TABS.VariableSet} tab="Variable Set">
            <ResponseVariableSet />
          </Tabs.TabPane>
          <Tabs.TabPane key={TABS.Metadata} tab="Metadata">
            <ResponseMetadata />
          </Tabs.TabPane>
        </Tabs>
      </S.TabsContainer>
    </S.Container>
  );
};

export default RunDetailTriggerData;
