import {Tabs} from 'antd';
import {useCallback, useState} from 'react';
import {useMount} from 'react-use';
import Drawer from 'components/Drawer';
import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import SetupAlert from 'components/SetupAlert';
import SpanDetail from 'components/SpanDetail';
import TestOutputs from 'components/TestOutputs';
import TestResults from 'components/TestResults';
import TestSpecDetail from 'components/TestSpecDetail';
import TestSpecForm from 'components/TestSpecForm';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import Switch from 'components/Visualization/components/Switch';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import TestRunAnalytics from 'services/Analytics/TestRunAnalytics.service';
import AssertionService from 'services/Assertion.service';
import {TAssertionResultEntry} from 'types/Assertion.types';
import {TTestRun} from 'types/TestRun.types';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import * as S from './RunDetailTest.styled';
import Visualization from './Visualization';

interface IProps {
  run: TTestRun;
  testId: string;
}

const RunDetailTest = ({run, testId}: IProps) => {
  const {selectedSpan, onSetFocusedSpan, onSelectSpan} = useSpan();
  const {remove, revert, selectedTestSpec, setSelectedSpec, setSelectorSuggestions, setPrevSelector} = useTestSpecs();
  const {isOpen: isTestSpecFormOpen, formProps, onSubmit, open, close} = useTestSpecForm();
  const [visualizationType, setVisualizationType] = useState(VisualizationType.Dag);
  const {
    state: {tourActive},
    setState,
  } = useGuidedTour();
  useMount(() => {
    if (tourActive) setState(st => ({...st, run: true, stepIndex: 3}));
  });

  const handleClose = useCallback(() => {
    onSetFocusedSpan('');
    setSelectedSpec();
  }, [onSetFocusedSpan, setSelectedSpec]);

  const handleEdit = useCallback(
    ({selector, resultList: list}: TAssertionResultEntry) => {
      AssertionAnalyticsService.onAssertionEdit();

      open({
        isEditing: true,
        selector,
        defaultValues: {
          assertions: list.map(({assertion}) => AssertionService.getStructuredAssertion(assertion)),
          selector,
        },
      });
    },
    [open]
  );

  const handleDelete = useCallback(
    (selector: string) => {
      AssertionAnalyticsService.onAssertionDelete();
      remove(selector);
    },
    [remove]
  );

  const handleRevert = useCallback(
    (originalSelector: string) => {
      AssertionAnalyticsService.onRevertAssertion();
      revert(originalSelector);
    },
    [revert]
  );

  const handleSelectSpan = useCallback(
    (spanId: string) => {
      onSelectSpan(spanId);
      onSetFocusedSpan(spanId);
    },
    [onSelectSpan, onSetFocusedSpan]
  );

  return (
    <S.Container>
      <SetupAlert />
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
                  onSubmit={values => {
                    setSelectorSuggestions([]);
                    setPrevSelector('');
                    onSubmit(values);
                  }}
                  runId={run.id}
                  testId={testId}
                  {...formProps}
                  onCancel={() => {
                    setSelectorSuggestions([]);
                    setPrevSelector('');
                    close();
                  }}
                  onClearSelectorSuggestions={() => {
                    setSelectorSuggestions([]);
                  }}
                  onClickPrevSelector={prevSelector => {
                    setPrevSelector(prevSelector);
                  }}
                />
              ) : (
                <S.TabsContainer>
                  <Tabs defaultActiveKey="testSpecs" size="small">
                    <Tabs.TabPane key="testSpecs" tab="Test Specs">
                      <TestResults onDelete={handleDelete} onEdit={handleEdit} onRevert={handleRevert} />
                    </Tabs.TabPane>
                    <Tabs.TabPane key="testOutputs" tab="Test Outputs">
                      <TestOutputs />
                    </Tabs.TabPane>
                  </Tabs>

                  <TestSpecDetail
                    isOpen={Boolean(selectedTestSpec)}
                    onClose={handleClose}
                    onDelete={handleDelete}
                    onEdit={handleEdit}
                    onRevert={handleRevert}
                    onSelectSpan={handleSelectSpan}
                    selectedSpan={selectedSpan?.id}
                    testSpec={selectedTestSpec}
                  />
                </S.TabsContainer>
              )}
            </S.SectionRight>
          </S.Container>
        }
      />
    </S.Container>
  );
};

export default RunDetailTest;
