import {useCallback, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useAppSelector} from 'redux/hooks';
import Drawer from 'components/Drawer';
import SpanDetail from 'components/SpanDetail';
import Switch from 'components/Visualization/components/Switch';
import {TestState} from 'constants/TestRun.constants';
import SpanSelectors from 'selectors/Span.selectors';
import TraceSelectors from 'selectors/Trace.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import TestRun from 'models/TestRun.model';
import * as S from './RunDetailTrace.styled';
import Search from './Search';
import Visualization from './Visualization';
import SetupAlert from '../SetupAlert';

interface IProps {
  run: TestRun;
  testId: string;
}

export enum VisualizationType {
  Dag,
  Timeline,
}

const RunDetailTrace = ({run, testId}: IProps) => {
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
          <S.Section>
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
              <Visualization runState={run.state} spans={run?.trace?.spans ?? []} type={visualizationType} />
            </S.VisualizationContainer>
          </S.Section>
        }
      />
    </S.Container>
  );
};

export default RunDetailTrace;
