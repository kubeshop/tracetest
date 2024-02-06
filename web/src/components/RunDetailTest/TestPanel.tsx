import {Tabs} from 'antd';
import {useCallback, useState} from 'react';
import {useSearchParams} from 'react-router-dom';
import {VisualizationType, getIsDAGDisabled} from 'components/RunDetailTrace/RunDetailTrace';
import TestOutputs from 'components/TestOutputs';
import TestOutputForm from 'components/TestOutputForm/TestOutputForm';
import TestResults from 'components/TestResults';
import TestSpecDetail from 'components/TestSpecDetail';
import TestSpecForm from 'components/TestSpecForm';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import Switch from 'components/Visualization/components/Switch';
import {TestState} from 'constants/TestRun.constants';
import {TAssertionResultEntry} from 'models/AssertionResults.model';
import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import TestRunAnalytics from 'services/Analytics/TestRunAnalytics.service';
import AssertionService from 'services/Assertion.service';
import {useTest} from 'providers/Test/Test.provider';
import * as S from './RunDetailTest.styled';
import Visualization from './Visualization';
import {FillPanel} from '../ResizablePanels';
import SkipTraceCollectionInfo from '../SkipTraceCollectionInfo';

const TABS = {
  SPECS: 'specs',
  OUTPUTS: 'outputs',
} as const;

interface IProps {
  run: TestRun;
  runEvents: TestRunEvent[];
  testId: string;
}

const TestPanel = ({run, testId, runEvents}: IProps) => {
  const [query, updateQuery] = useSearchParams();
  const {selectedSpan, onSetFocusedSpan, onSelectSpan} = useSpan();
  const {remove, revert, selectedTestSpec, setSelectedSpec, setSelectorSuggestions, setPrevSelector, specs} =
    useTestSpecs();
  const {isOpen: isTestSpecFormOpen, formProps, onSubmit, open, close, isValid, onIsValid} = useTestSpecForm();
  const {
    isEditing,
    isLoading,
    isOpen: isTestOutputFormOpen,
    onClose,
    onSubmit: onSubmitTestOutput,
    output,
    outputs,
    isValid: isOutputFormValid,
    onValidate,
  } = useTestOutput();
  const {
    test: {skipTraceCollection},
  } = useTest();

  const isDAGDisabled = getIsDAGDisabled(run?.trace?.spans?.length);
  const [visualizationType, setVisualizationType] = useState(() =>
    isDAGDisabled ? VisualizationType.Timeline : VisualizationType.Dag
  );

  const handleClose = useCallback(() => {
    onSetFocusedSpan('');
    setSelectedSpec();
  }, [onSetFocusedSpan, setSelectedSpec]);

  const handleEdit = useCallback(
    ({selector, resultList: list}: TAssertionResultEntry, name: string) => {
      AssertionAnalyticsService.onAssertionEdit();

      open({
        isEditing: true,
        selector,
        defaultValues: {
          assertions: list.map(({assertion}) => AssertionService.getStructuredAssertion(assertion)),
          selector,
          name,
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
    <FillPanel>
      <S.Container>
        <S.SectionLeft>
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
            spans={run?.trace?.spans ?? []}
            type={visualizationType}
          />
        </S.SectionLeft>

        <S.SectionRight>
          {isTestSpecFormOpen && (
            <TestSpecForm
              onSubmit={values => {
                setSelectorSuggestions([]);
                setPrevSelector('');
                onSubmit(values);
              }}
              isValid={isValid}
              onIsValid={onIsValid}
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
              isValid={isOutputFormValid}
              onValidate={onValidate}
            />
          )}

          {!isTestSpecFormOpen && !isTestOutputFormOpen && (
            <S.TabsContainer>
              <Tabs
                activeKey={query.get('tab') || TABS.SPECS}
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
                      Test Outputs <S.CountBadge data-cy="output-count-badge" count={outputs.length} />
                    </>
                  }
                >
                  <TestOutputs outputs={outputs} />
                </Tabs.TabPane>
              </Tabs>
            </S.TabsContainer>
          )}

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
        </S.SectionRight>
      </S.Container>
    </FillPanel>
  );
};

export default TestPanel;
