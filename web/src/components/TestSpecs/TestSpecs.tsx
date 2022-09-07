import TestSpec from 'components/TestSpec';
import {TAssertionResultEntry, TAssertionResults} from 'types/Assertion.types';
import Empty from './Empty';
import * as S from './TestSpecs.styled';

interface IProps {
  assertionResults: TAssertionResults;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
  onOpen(selector: string): void;
  onRevert(originalSelector: string): void;
}

const TestSpecs = ({assertionResults: {resultList}, onDelete, onEdit, onOpen, onRevert}: IProps) => {
  if (!resultList.length) {
    return <Empty />;
  }

  return (
    <S.Container data-cy="test-specs-container">
      {resultList.map(specResult =>
        specResult.resultList.length ? (
          <TestSpec
            key={specResult.id}
            onDelete={onDelete}
            onEdit={onEdit}
            onOpen={onOpen}
            onRevert={onRevert}
            testSpec={specResult}
          />
        ) : null
      )}
    </S.Container>
  );
};

export default TestSpecs;
