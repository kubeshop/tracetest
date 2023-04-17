import {useMemo} from 'react';

import {useAppSelector} from 'redux/hooks';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import AssertionService from 'services/Assertion.service';
import {TAssertionResultEntry} from 'models/AssertionResults.model';
import Actions from './Actions';
import Header from './Header';
import * as S from './TestSpec.styled';

interface IProps {
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry, name: string): void;
  onOpen(selector: string): void;
  onRevert(originalSelector: string): void;
  testSpec: TAssertionResultEntry;
}

const TestSpec = ({
  onDelete,
  onEdit,
  onOpen,
  onRevert,
  testSpec,
  testSpec: {resultList, selector, spanIds},
}: IProps) => {
  const {
    isDeleted = false,
    isDraft = false,
    originalSelector = '',
    name = '',
  } = useAppSelector(state => TestSpecsSelectors.selectSpecBySelector(state, selector)) || {};
  const totalPassedChecks = useMemo(() => AssertionService.getTotalPassedChecks(resultList), [resultList]);
  const hasError = useMemo(() => AssertionService.hasError(resultList), [resultList]);

  return (
    <S.Container
      $isDeleted={isDeleted}
      data-cy="test-spec-container"
      onClick={() => {
        onOpen(selector);
      }}
    >
      <Header
        affectedSpans={spanIds.length}
        assertionsFailed={totalPassedChecks?.false ?? 0}
        assertionsPassed={totalPassedChecks?.true ?? 0}
        hasError={hasError}
        selector={selector}
        title={!selector && !name ? 'All Spans' : name}
      />
      <Actions
        isDeleted={isDeleted}
        isDraft={isDraft}
        onDelete={() => {
          onDelete(selector);
        }}
        onEdit={() => {
          onEdit(testSpec, name);
        }}
        onRevert={() => {
          onRevert(originalSelector);
        }}
      />
    </S.Container>
  );
};

export default TestSpec;
