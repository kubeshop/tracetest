import {useState} from 'react';
import TestRun from 'models/TestRun.model';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {TestState} from 'constants/TestRun.constants';
import TestRunEvent from 'models/TestRunEvent.model';
import Search from './Search';
import {VisualizationType, getIsDAGDisabled} from './RunDetailTrace';
import * as S from './RunDetailTrace.styled';
import Switch from '../Visualization/components/Switch/Switch';
import Visualization from './Visualization';
import {FillPanel} from '../ResizablePanels';
import SkipTraceCollectionInfo from '../SkipTraceCollectionInfo';

type TProps = {
  run: TestRun;
  testId: string;
  runEvents: TestRunEvent[];
  skipTraceCollection: boolean;
};

const TracePanel = ({run, testId, runEvents, skipTraceCollection}: TProps) => {
  const isDAGDisabled = getIsDAGDisabled(run?.trace?.spans?.length);
  const [visualizationType, setVisualizationType] = useState(() =>
    isDAGDisabled ? VisualizationType.Timeline : VisualizationType.Dag
  );

  return (
    <FillPanel>
      <S.Container>
        <S.SectionLeft $hasShadow>
          <S.SearchContainer>
            <Search runId={run.id} testId={testId} />
          </S.SearchContainer>

          <S.VisualizationContainer>
            {skipTraceCollection && <SkipTraceCollectionInfo runId={run.id} testId={testId} />}
            <S.SwitchContainer>
              {run.state === TestState.FINISHED && (
                <Switch
                  isDAGDisabled={isDAGDisabled}
                  onChange={type => {
                    TraceAnalyticsService.onSwitchDiagramView(type);
                    setVisualizationType(type);
                  }}
                  type={visualizationType}
                  totalSpans={run?.trace?.spans?.length}
                />
              )}
            </S.SwitchContainer>
            <Visualization
              isDAGDisabled={isDAGDisabled}
              runEvents={runEvents}
              runState={run.state}
              spans={run?.trace?.spans ?? []}
              type={visualizationType}
            />
          </S.VisualizationContainer>
        </S.SectionLeft>
      </S.Container>
    </FillPanel>
  );
};

export default TracePanel;
