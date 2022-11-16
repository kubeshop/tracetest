import {Tabs} from 'antd';
import {TriggerTypes} from 'constants/Test.constants';
import {TestState} from 'constants/TestRun.constants';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {TTriggerResult} from 'types/Test.types';
import {TTestRunState} from 'types/TestRun.types';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import {Steps} from '../GuidedTour/traceStepList';
import ResponseBody from './ResponseBody';
import ResponseHeaders from './ResponseHeaders';
import ResponseOutputs from './ResponseOutputs';
import * as S from './RunDetailTriggerResponse.styled';

const isTransactionsEnabled = ExperimentalFeature.isEnabled('transactions');

interface IProps {
  state: TTestRunState;
  triggerResult?: TTriggerResult;
  triggerTime?: number;
}

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
          defaultActiveKey="1"
          data-cy="run-detail-trigger-response"
          size="small"
          onChange={newTab => TestRunAnalyticsService.onTriggerResponseTabChange(newTab)}
        >
          <Tabs.TabPane key="1" tab="Body">
            <ResponseBody body={body} bodyMimeType={bodyMimeType} />
          </Tabs.TabPane>
          <Tabs.TabPane key="2" tab="Headers">
            <ResponseHeaders headers={headers} />
          </Tabs.TabPane>
          {isTransactionsEnabled && (
            <Tabs.TabPane key="3" tab="Outputs">
              <ResponseOutputs />
            </Tabs.TabPane>
          )}
        </Tabs>
      </S.TabsContainer>
    </S.Container>
  );
};

export default RunDetailTriggerResponse;
