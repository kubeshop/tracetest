import {Tabs} from 'antd';
import {useSearchParams} from 'react-router-dom';
import {StepsID} from 'components/GuidedTour/testRunSteps';
import {TriggerTypes} from 'constants/Test.constants';
import {TestState} from 'constants/TestRun.constants';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import ResponseBody from './ResponseBody';
import ResponseEnvironment from './ResponseEnvironment';
import ResponseHeaders from './ResponseHeaders';
import * as S from './RunDetailTriggerResponse.styled';
import {IPropsComponent} from './RunDetailTriggerResponseFactory';

const TabsKeys = {
  Body: 'body',
  Headers: 'headers',
  Environment: 'environment',
};

const RunDetailTriggerResponse = ({
  state,
  triggerTime = 0,
  triggerResult: {headers, body = '', statusCode = 200, bodyMimeType} = {
    body: '',
    type: TriggerTypes.http,
    statusCode: 200,
    bodyMimeType: '',
  },
}: IPropsComponent) => {
  const [query, updateQuery] = useSearchParams();

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>Response Data</S.Title>
        <div>
          <S.StatusText>
            Status: <S.StatusSpan $isError={statusCode >= 400}>{statusCode}</S.StatusSpan>
          </S.StatusText>
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
          defaultActiveKey={query.get('tab') || TabsKeys.Body}
          data-cy="run-detail-trigger-response"
          size="small"
          onChange={newTab => {
            TestRunAnalyticsService.onTriggerResponseTabChange(newTab);
            updateQuery([['tab', newTab]]);
          }}
        >
          <Tabs.TabPane key={TabsKeys.Body} tab="Body">
            <ResponseBody body={body} bodyMimeType={bodyMimeType} state={state} />
          </Tabs.TabPane>
          <Tabs.TabPane key={TabsKeys.Headers} tab="Headers">
            <ResponseHeaders headers={headers} state={state} />
          </Tabs.TabPane>
          <Tabs.TabPane key={TabsKeys.Environment} tab="Environment">
            <ResponseEnvironment />
          </Tabs.TabPane>
        </Tabs>
      </S.TabsContainer>
    </S.Container>
  );
};

export default RunDetailTriggerResponse;
