import {useState} from 'react';
import {VisualizationType, getIsDAGDisabled} from 'components/RunDetailTrace/RunDetailTrace';
import Switch from 'components/Visualization/components/Switch';
import {TestState} from 'constants/TestRun.constants';
import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import TestRunAnalytics from 'services/Analytics/TestRunAnalytics.service';
import {useTest} from 'providers/Test/Test.provider';
import * as S from './RunDetailTest.styled';
import Visualization from './Visualization';
import {FillPanel} from '../ResizablePanels';
import SkipTraceCollectionInfo from '../SkipTraceCollectionInfo';

interface IProps {
  run: TestRun;
  runEvents: TestRunEvent[];
  testId: string;
}

const TestPanel = ({run, testId, runEvents}: IProps) => {
  const isDAGDisabled = getIsDAGDisabled(run?.trace?.spans?.length);
  const [visualizationType, setVisualizationType] = useState(() =>
    isDAGDisabled ? VisualizationType.Timeline : VisualizationType.Dag
  );

  const {
    test: {skipTraceCollection},
  } = useTest();

  return (
    <FillPanel>
      <S.Container>
        <S.SectionLeft $isTimeline={visualizationType === VisualizationType.Timeline}>
          <S.SwitchContainer>
            {run.state === TestState.FINISHED && (
              <Switch
                isDAGDisabled={isDAGDisabled}
                onChange={type => {
                  TestRunAnalytics.onSwitchDiagramView(type);
                  setVisualizationType(type);
                }}
                type={visualizationType}
                totalSpans={run?.trace?.spans?.length}
              />
            )}
          </S.SwitchContainer>

          {skipTraceCollection && <SkipTraceCollectionInfo runId={run.id} testId={testId} />}
          <Visualization
            isDAGDisabled={isDAGDisabled}
            runEvents={runEvents}
            runState={run.state}
            trace={run.trace}
            type={visualizationType}
          />
        </S.SectionLeft>
      </S.Container>
    </FillPanel>
  );
};

export default TestPanel;
