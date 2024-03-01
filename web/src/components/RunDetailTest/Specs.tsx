import {Tabs} from 'antd';
import TestOutputs from 'components/TestOutputs';
import TestOutputForm from 'components/TestOutputForm/TestOutputForm';
import TestResults from 'components/TestResults';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {TAssertionResultEntry} from 'models/AssertionResults.model';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import AssertionService from 'services/Assertion.service';
import {useSearchParams} from 'react-router-dom';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import TestSpecDetail from 'components/TestSpecDetail';
import {useTest} from 'providers/Test/Test.provider';
import TestRun from 'models/TestRun.model';
import {useSpan} from 'providers/Span/Span.provider';
import {useCallback} from 'react';
import TestSpecForm from 'components/TestSpecForm';
import * as S from './RunDetailTest.styled';
import SpanResultDetail from '../SpanResultDetail';
import {useTestSpecForm} from '../TestSpecForm/TestSpecForm.provider';

const TABS = {
  SPECS: 'specs',
  OUTPUTS: 'outputs',
} as const;

interface IProps {
  run: TestRun;
}

const Specs = ({run}: IProps) => {
  const [query, updateQuery] = useSearchParams();
  const {selectedSpan, onSetFocusedSpan} = useSpan();
  const {
    remove,
    revert,
    selectedTestSpec,
    selectedSpanResult,
    setSelectedSpanResult,
    setSelectedSpec,
    setSelectorSuggestions,
    setPrevSelector,
    specs,
  } = useTestSpecs();
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
    test: {id: testId},
  } = useTest();

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

  return (
    <>
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
        selectedSpan={selectedSpan?.id}
        testSpec={selectedTestSpec}
      />

      <SpanResultDetail
        isOpen={!!selectedSpanResult}
        spanResult={selectedSpanResult}
        onClose={() => setSelectedSpanResult()}
      />
    </>
  );
};

export default Specs;
