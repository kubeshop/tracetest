import {Tabs} from 'antd';
import {useSearchParams} from 'react-router-dom';
import {TriggerTypes} from 'constants/Test.constants';
import {TestState} from 'constants/TestRun.constants';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {TTriggerResult} from 'types/Test.types';
import {TTestRunState} from 'types/TestRun.types';
import {Steps} from '../GuidedTour/traceStepList';
import ResponseBody from './ResponseBody';
import ResponseEnvironment from './ResponseEnvironment';
import ResponseHeaders from './ResponseHeaders';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  state: TTestRunState;
  triggerResult?: TTriggerResult;
  triggerTime?: number;
}

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
}: IProps) => {
  const [query, updateQuery] = useSearchParams();

  return (
    <S.Container data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Graph)}>
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
      <S.TabsContainer>
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
            <ResponseBody body={body} bodyMimeType={bodyMimeType} />
          </Tabs.TabPane>
          <Tabs.TabPane key={TabsKeys.Headers} tab="Headers">
            <ResponseHeaders headers={headers} />
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
