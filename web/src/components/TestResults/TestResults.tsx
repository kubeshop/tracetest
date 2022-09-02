import {useCallback} from 'react';

import LoadingSpinner from 'components/LoadingSpinner';
import TestSpecDetail from 'components/TestSpecDetail';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import TestSpecs from 'components/TestSpecs';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import {TAssertionResultEntry} from 'types/Assertion.types';
import Header from './Header';
import * as S from './TestResults.styled';

const TestResults = () => {
  const {open} = useTestSpecForm();
  const {isLoading, assertionResults, remove, revert, setSelectedSpec} = useTestSpecs();
  const {selectedSpan, onSetFocusedSpan} = useSpan();
  const {totalFailedSpecs, totalPassedSpecs} = useAppSelector(TestDefinitionSelectors.selectTotalSpecs);
  const selectedAssertion = useAppSelector(TestDefinitionSelectors.selectSelectedAssertion);
  const selectedTestSpec = useAppSelector(state =>
    TestDefinitionSelectors.selectAssertionBySelector(state, selectedAssertion ?? '')
  );

  const handleOpen = useCallback(
    (selector: string) => {
      AssertionAnalyticsService.onAssertionClick();
      const testSpec = assertionResults?.resultList?.find(specResult => specResult.selector === selector);
      onSetFocusedSpan('');
      setSelectedSpec(testSpec);
    },
    [assertionResults?.resultList, onSetFocusedSpan, setSelectedSpec]
  );

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
          assertions: list.map(({assertion}) => assertion),
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
      onSetFocusedSpan(spanId);
    },
    [onSetFocusedSpan]
  );

  return (
    <S.Container>
      <Header selectedSpan={selectedSpan!} totalFailedSpecs={totalFailedSpecs} totalPassedSpecs={totalPassedSpecs} />

      {(isLoading || !assertionResults) && (
        <S.LoadingContainer>
          <LoadingSpinner />
        </S.LoadingContainer>
      )}

      {!isLoading && Boolean(assertionResults) && (
        <TestSpecs
          assertionResults={assertionResults!}
          onDelete={handleDelete}
          onEdit={handleEdit}
          onOpen={handleOpen}
          onRevert={handleRevert}
        />
      )}

      <TestSpecDetail
        isOpen={Boolean(selectedAssertion)}
        onClose={handleClose}
        onDelete={handleDelete}
        onEdit={handleEdit}
        onRevert={handleRevert}
        onSelectSpan={handleSelectSpan}
        selectedSpan={selectedSpan?.id}
        testSpec={selectedTestSpec}
      />
    </S.Container>
  );
};

export default TestResults;
