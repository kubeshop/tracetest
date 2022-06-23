import {Typography} from 'antd';
import {useCallback} from 'react';
import {useTestDefinition} from '../../providers/TestDefinition/TestDefinition.provider';
import AssertionAnalyticsService from '../../services/Analytics/AssertionAnalytics.service';
import {TAssertionResultEntry, TAssertionResults} from '../../types/Assertion.types';
import AssertionCard from '../AssertionCard/AssertionCard';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {OPEN_BOTTOM_PANEL_STATE, useRunLayout} from '../RunLayout';
import * as S from './AssertionCardList.styled';

interface TAssertionCardListProps {
  assertionResults: TAssertionResults;
  onSelectSpan(spanId: string): void;
  testId: string;
}

const AssertionCardList: React.FC<TAssertionCardListProps> = ({assertionResults: {resultList}, onSelectSpan}) => {
  const {open} = useAssertionForm();
  const {remove} = useTestDefinition();
  const {openBottomPanel} = useRunLayout();

  const handleEdit = useCallback(
    ({selector, resultList: list, selectorList, pseudoSelector}: TAssertionResultEntry) => {
      AssertionAnalyticsService.onAssertionEdit();
      openBottomPanel(OPEN_BOTTOM_PANEL_STATE.FORM);
      open({
        isEditing: true,
        selector,
        defaultValues: {
          assertionList: list.map(({assertion}) => assertion),
          selectorList,
          pseudoSelector,
        },
      });
    },
    [open, openBottomPanel]
  );

  const handleDelete = useCallback(
    (selector: string) => {
      AssertionAnalyticsService.onAssertionDelete();
      remove(selector);
    },
    [remove]
  );

  return (
    <S.AssertionCardList data-cy="assertion-card-list">
      {resultList.length ? (
        resultList.map(assertionResult =>
          assertionResult.resultList.length ? (
            <AssertionCard
              key={assertionResult.id}
              assertionResult={assertionResult}
              onSelectSpan={onSelectSpan}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          ) : null
        )
      ) : (
        <S.EmptyStateContainer data-cy="empty-assertion-card-list">
          <S.EmptyStateIcon />
          <Typography.Text disabled>Add an Assertion to See the Result</Typography.Text>
        </S.EmptyStateContainer>
      )}
    </S.AssertionCardList>
  );
};

export default AssertionCardList;
