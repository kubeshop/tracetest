import {Tabs} from 'antd';
import {useCallback, useState} from 'react';
import {useSearchParams} from 'react-router-dom';
import {useMount} from 'react-use';

import Drawer from 'components/Drawer';
import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import SetupAlert from 'components/SetupAlert';
import SpanDetail from 'components/SpanDetail';
import TestOutputs from 'components/TestOutputs';
import TestOutputForm from 'components/TestOutputForm/TestOutputForm';
import TestResults from 'components/TestResults';
import TestSpecDetail from 'components/TestSpecDetail';
import TestSpecForm from 'components/TestSpecForm';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import Switch from 'components/Visualization/components/Switch';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import TestRunAnalytics from 'services/Analytics/TestRunAnalytics.service';
import AssertionService from 'services/Assertion.service';
import {TAssertionResultEntry} from 'types/Assertion.types';
import {TTestRun} from 'types/TestRun.types';
import * as S from './RunDetailTest.styled';
import Visualization from './Visualization';

const TABS = {
  SPECS: 'specs',
  OUTPUTS: 'outputs',
} as const;

interface IProps {
  run: TTestRun;
  testId: string;
}

const RunDetailTest = ({run, testId}: IProps) => {
  const {selectedSpan, onSetFocusedSpan, onSelectSpan} = useSpan();
  const {remove, revert, selectedTestSpec, setSelectedSpec, setSelectorSuggestions, setPrevSelector, specs} =
    useTestSpecs();
  const {isOpen: isTestSpecFormOpen, formProps, onSubmit, open, close} = useTestSpecForm();
  const {
    isEditing,
    isLoading,
    isOpen: isTestOutputFormOpen,
    onClose,
    onSubmit: onSubmitTestOutput,
    output,
    outputs,
  } = useTestOutput();
  const [visualizationType, setVisualizationType] = useState(VisualizationType.Dag);
  const {
    state: {tourActive},
    setState,
  } = useGuidedTour();
  useMount(() => {
    if (tourActive) setState(st => ({...st, run: true, stepIndex: 3}));
  });
  const [query, updateQuery] = useSearchParams();

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
              {isTestSpecFormOpen && (
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
              )}

              {isTestOutputFormOpen && (
                <TestOutputForm
                  isEditing={isEditing}
                  isLoading={isLoading}
                  onCancel={onClose}
                  onSubmit={onSubmitTestOutput}
                  output={output}
                  runId={run.id}
                  testId={testId}
                />
              )}

              {!isTestSpecFormOpen && !isTestOutputFormOpen && (
                <S.TabsContainer>
                  <Tabs
                    defaultActiveKey={query.get('tab') || TABS.SPECS}
                    onChange={tab =>
                      updateQuery(
                        selectedSpan
                          ? [
                              ['selectedSpan', selectedSpan.id],
                              ['tab', tab],
                            ]
                          : [['tab', tab]]
                      )
                    }
                    size="small"
                  >
                    <Tabs.TabPane
                      key={TABS.SPECS}
                      tab={
                        <>
                          Test Specs <S.CountBadge count={specs.length} />
                        </>
                      }
                    >
                      <TestResults onDelete={handleDelete} onEdit={handleEdit} onRevert={handleRevert} />
                    </Tabs.TabPane>
                    <Tabs.TabPane
                      key={TABS.OUTPUTS}
                      tab={
                        <>
                          Test Outputs <S.CountBadge count={outputs.length} />
                        </>
                      }
                    >
                      <TestOutputs outputs={outputs} />
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
