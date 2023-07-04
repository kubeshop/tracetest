import {useState} from 'react';
import TestRun from 'models/TestRun.model';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {TestState} from 'constants/TestRun.constants';
import TestRunEvent from 'models/TestRunEvent.model';
import Search from './Search';
import {VisualizationType} from './RunDetailTrace';
import * as S from './RunDetailTrace.styled';
import Switch from '../Visualization/components/Switch/Switch';
import Visualization from './Visualization';
import {FillPanel} from '../ResizablePanels';

type TProps = {
  run: TestRun;
  testId: string;
  runEvents: TestRunEvent[];
};

const TracePanel = ({run, testId, runEvents}: TProps) => {
  const [visualizationType, setVisualizationType] = useState(VisualizationType.Dag);

  return (
    <FillPanel>
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
      </S.Container>
    </FillPanel>
  );
};

export default TracePanel;
