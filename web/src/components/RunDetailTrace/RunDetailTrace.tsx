import {useCallback, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useAppSelector} from 'redux/hooks';
import Drawer from 'components/Drawer';
import SpanDetail from 'components/SpanDetail';
import SkeletonResponse from 'components/RunDetailTriggerResponse/SkeletonResponse';
import Switch from 'components/Visualization/components/Switch';
import {TestState} from 'constants/TestRun.constants';
import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import Trace from 'models/Trace.model';
import TestRunEvent from 'models/TestRunEvent.model';
import SpanSelectors from 'selectors/Span.selectors';
import TraceSelectors from 'selectors/Trace.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import LintResults from './LintResults';
import * as S from './RunDetailTrace.styled';
import Search from './Search';
import Visualization from './Visualization';
import SetupAlert from '../SetupAlert';

interface IProps {
  run: TestRun;
  runEvents: TestRunEvent[];
  testId: string;
}

export enum VisualizationType {
  Dag,
  Timeline,
}

const RunDetailTrace = ({run, runEvents, testId}: IProps) => {
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);
  const searchText = useAppSelector(TraceSelectors.selectSearchText);
  const span = useAppSelector(state => SpanSelectors.selectSpanById(state, selectedSpan, testId, run.id));
  const navigate = useNavigate();
  const [visualizationType, setVisualizationType] = useState(VisualizationType.Dag);

  const handleOnCreateSpec = useCallback(() => {
    navigate(`/test/${testId}/run/${run.id}/test`);
  }, [navigate, run.id, testId]);

  return (
    <S.Container>
      <SetupAlert />
      <Drawer
        leftPanel={<SpanDetail onCreateTestSpec={handleOnCreateSpec} searchText={searchText} span={span} />}
        rightPanel={
          <S.Container>
            <S.SectionLeft>
              <S.SearchContainer>
                <Search runId={run.id} testId={testId} />
              </S.SearchContainer>

              <S.VisualizationContainer>
                <S.SwitchContainer>
                  {run.state === TestState.FINISHED && (
                    <Switch
                      onChange={type => {
                        TraceAnalyticsService.onSwitchDiagramView(type);
                        setVisualizationType(type);
                      }}
                      type={visualizationType}
                    />
                  )}
                </S.SwitchContainer>
                <Visualization
                  runEvents={runEvents}
                  runState={run.state}
                  spans={run?.trace?.spans ?? []}
                  type={visualizationType}
                />
              </S.VisualizationContainer>
            </S.SectionLeft>

            <S.SectionRight $shouldScroll>
              {isRunStateFinished(run.state) ? (
                <LintResults linterResult={run.linter} trace={run?.trace ?? Trace({})} />
              ) : (
                <SkeletonResponse />
              )}
            </S.SectionRight>
          </S.Container>
        }
      />
    </S.Container>
  );
};

export default RunDetailTrace;
