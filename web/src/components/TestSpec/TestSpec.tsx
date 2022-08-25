import {useMemo} from 'react';

import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import AssertionService from 'services/Assertion.service';
import {TAssertionResultEntry} from 'types/Assertion.types';
import Actions from './Actions';
import Header from './Header';
import * as S from './TestSpec.styled';

interface IProps {
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
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
  } = useAppSelector(state => TestDefinitionSelectors.selectDefinitionBySelector(state, selector)) || {};
  const totalPassedChecks = useMemo(() => AssertionService.getTotalPassedChecks(resultList), [resultList]);

  return (
    <S.Container
      data-cy="assertion-card"
      onClick={() => {
        onOpen(selector);
      }}
    >
      <Header
        affectedSpans={spanIds.length}
        assertionsFailed={totalPassedChecks?.false ?? 0}
        assertionsPassed={totalPassedChecks?.true ?? 0}
        title={selector}
      />
      <Actions
        isDeleted={isDeleted}
        isDraft={isDraft}
        onDelete={() => {
          onDelete(selector);
        }}
        onEdit={() => {
          onEdit(testSpec);
        }}
        onRevert={() => {
          onRevert(originalSelector);
        }}
      />
    </S.Container>
  );
};

export default TestSpec;
