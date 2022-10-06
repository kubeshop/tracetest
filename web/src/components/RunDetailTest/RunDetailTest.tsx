import {useState} from 'react';
import {useMount} from 'react-use';
import Drawer from 'components/Drawer';
import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import SpanDetail from 'components/SpanDetail';
import TestResults from 'components/TestResults';
import TestSpecForm from 'components/TestSpecForm';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import Switch from 'components/Visualization/components/Switch';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import TestRunAnalytics from 'services/Analytics/TestRunAnalytics.service';
import {TTestRun} from 'types/TestRun.types';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import * as S from './RunDetailTest.styled';
import Visualization from './Visualization';

interface IProps {
  run: TTestRun;
  testId: string;
}

const RunDetailTest = ({run, testId}: IProps) => {
  const {selectedSpan} = useSpan();
  const {selectedTestSpec} = useTestSpecs();
  const {isOpen: isTestSpecFormOpen, formProps, onSubmit, close} = useTestSpecForm();
  const [visualizationType, setVisualizationType] = useState(VisualizationType.Dag);
  const {
    state: {tourActive},
    setState,
  } = useGuidedTour();
  useMount(() => {
    if (tourActive) setState(st => ({...st, run: true, stepIndex: 3}));
  });
  return (
    <S.Container>
      <Drawer
        leftPanel={<SpanDetail span={selectedSpan} />}
        rightPanel={
          <S.Container>
            <S.SectionLeft>
              <S.SwitchContainer>
                <Switch
                  onChange={type => {
                    TestRunAnalytics.onSwitchDiagramView(type);
                    setVisualizationType(type);
                  }}
                  type={visualizationType}
                />
              </S.SwitchContainer>

              <Visualization runState={run.state} spans={run?.trace?.spans ?? []} type={visualizationType} />
            </S.SectionLeft>

            <S.SectionRight $shouldScroll={!selectedTestSpec}>
              {isTestSpecFormOpen ? (
                <TestSpecForm
                  onSubmit={onSubmit}
                  runId={run.id}
                  testId={testId}
                  {...formProps}
                  onCancel={() => {
                    close();
                  }}
                />
              ) : (
                <TestResults />
              )}
            </S.SectionRight>
          </S.Container>
        }
      />
    </S.Container>
  );
};

export default RunDetailTest;
