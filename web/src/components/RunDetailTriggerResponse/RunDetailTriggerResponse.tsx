import {Tabs} from 'antd';
import {useCallback} from 'react';
import {useNavigate, useSearchParams} from 'react-router-dom';

import AttributeActions from 'components/AttributeActions';
import {StepsID} from 'components/GuidedTour/testRunSteps';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {CompareOperatorSymbolMap} from 'constants/Operator.constants';
import {TriggerTypes} from 'constants/Test.constants';
import {TestState} from 'constants/TestRun.constants';
import TestOutput from 'models/TestOutput.model';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import AssertionService from 'services/Assertion.service';
import {TSpanFlatAttribute} from 'types/Span.types';
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

const tracetestTriggerSelector = 'span[tracetest.span.type="general" name="Tracetest trigger"]';

const RunDetailTriggerResponse = ({
  runId,
  state,
  testId,
  triggerTime = 0,
  triggerResult: {headers, body = '', statusCode = 200, bodyMimeType} = {
    body: '',
    type: TriggerTypes.http,
    statusCode: 200,
    bodyMimeType: '',
  },
}: IPropsComponent) => {
  const navigate = useNavigate();
  const [query, updateQuery] = useSearchParams();
  const {onNavigateAndOpen} = useTestOutput();
  const {open} = useTestSpecForm();

  const handleCreateTestOutput = useCallback(
    ({key}: TSpanFlatAttribute) => {
      TestRunAnalyticsService.onAddAssertionButtonClick();
      const selector = tracetestTriggerSelector;

      const output = TestOutput({
        selector,
        selectorParsed: {query: selector},
        name: key,
        value: `attr:${key}`,
      });

      onNavigateAndOpen({...output});
    },
    [onNavigateAndOpen]
  );

  const handleCreateTestSpec = useCallback(
    ({value, key}: TSpanFlatAttribute) => {
      TestRunAnalyticsService.onAddAssertionButtonClick();
      const selector = tracetestTriggerSelector;

      open({
        isEditing: false,
        selector,
        defaultValues: {
          assertions: [
            {
              left: `attr:${key}`,
              comparator: CompareOperatorSymbolMap.EQUALS,
              right: AssertionService.extractExpectedString(value) || '',
            },
          ],
          selector,
        },
      });

      navigate(`/test/${testId}/run/${runId}/test`);
    },
    [navigate, open, runId, testId]
  );

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>Response Data</S.Title>
        <div>
          <AttributeActions
            attribute={{key: 'tracetest.response.status', value: `${statusCode}`}}
            onCreateTestOutput={handleCreateTestOutput}
            onCreateTestSpec={handleCreateTestSpec}
          >
            <S.StatusText>
              Status: <S.StatusSpan $isError={statusCode >= 400}>{statusCode}</S.StatusSpan>
            </S.StatusText>
          </AttributeActions>
          <AttributeActions
            attribute={{key: 'tracetest.span.duration', value: `${triggerTime}ms`}}
            onCreateTestOutput={handleCreateTestOutput}
            onCreateTestSpec={handleCreateTestSpec}
          >
            <S.StatusText>
              Time:{' '}
              <S.StatusSpan $isError={triggerTime > 1000}>
                {state === TestState.CREATED || state === TestState.EXECUTING ? '-' : `${triggerTime}ms`}
              </S.StatusSpan>
            </S.StatusText>
          </AttributeActions>
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
            <ResponseBody
              body={body}
              bodyMimeType={bodyMimeType}
              state={state}
              onCreateTestOutput={handleCreateTestOutput}
              onCreateTestSpec={handleCreateTestSpec}
            />
          </Tabs.TabPane>
          <Tabs.TabPane key={TabsKeys.Headers} tab="Headers">
            <ResponseHeaders
              headers={headers}
              state={state}
              onCreateTestOutput={handleCreateTestOutput}
              onCreateTestSpec={handleCreateTestSpec}
            />
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
