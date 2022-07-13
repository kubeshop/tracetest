import {Typography} from 'antd';
import {useCallback} from 'react';

import AssertionItem from 'components/AssertionItem';
import {useAssertionForm} from 'components/AssertionForm/AssertionForm.provider';
import {OPEN_BOTTOM_PANEL_STATE, useRunLayout} from 'components/RunLayout';
import {ResultViewModes} from 'constants/Test.constants';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import {TAssertionResultEntry, TAssertionResults} from 'types/Assertion.types';
import * as S from './AssertionList.styled';

interface IProps {
  assertionResults: TAssertionResults;
  onSelectSpan(spanId: string): void;
}

const AssertionList = ({assertionResults: {resultList}, onSelectSpan}: IProps) => {
  const selectedAssertion = useAppSelector(TestDefinitionSelectors.selectSelectedAssertion);
  const {open} = useAssertionForm();
  const {openBottomPanel} = useRunLayout();
  const {onSetFocusedSpan, selectedSpan} = useSpan();
  const {setSelectedAssertion, revert, remove, viewResultsMode} = useTestDefinition();
  const {
    run: {trace},
  } = useTestRun();

  const handleEdit = useCallback(
    ({selector, resultList: list, selectorList, pseudoSelector, isAdvancedSelector}: TAssertionResultEntry) => {
      AssertionAnalyticsService.onAssertionEdit();
      openBottomPanel(OPEN_BOTTOM_PANEL_STATE.FORM);

      open({
        isEditing: true,
        selector,
        defaultValues: {
          assertionList: list.map(({assertion}) => assertion),
          selectorList,
          selector,
          pseudoSelector,
          isAdvancedSelector: viewResultsMode === ResultViewModes.Advanced || isAdvancedSelector,
        },
      });
    },
    [open, openBottomPanel, viewResultsMode]
  );

  const handleDelete = useCallback(
    (selector: string) => {
      AssertionAnalyticsService.onAssertionDelete();
      remove(selector);
    },
    [remove]
  );

  return (
    <S.Container data-cy="assertion-card-list">
      {resultList.length ? (
        resultList.map(assertionResult =>
          assertionResult.resultList.length ? (
            <AssertionItem
              assertionResult={assertionResult}
              key={assertionResult.id}
              onDelete={handleDelete}
              onEdit={handleEdit}
              onRevertAssertion={revert}
              onSelectSpan={onSelectSpan}
              onSetFocusedSpan={onSetFocusedSpan}
              onSetSelectedAssertion={setSelectedAssertion}
              selectedAssertion={selectedAssertion}
              selectedSpan={selectedSpan?.id}
              trace={trace}
            />
          ) : null
        )
      ) : (
        <S.EmptyStateContainer data-cy="empty-assertion-card-list">
          <S.EmptyStateIcon />
          <Typography.Text disabled>Add an Assertion to See the Result</Typography.Text>
        </S.EmptyStateContainer>
      )}
    </S.Container>
  );
};

export default AssertionList;
