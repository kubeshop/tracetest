import TestSpec from 'components/TestSpec';
import AssertionResults, {TAssertionResultEntry} from 'models/AssertionResults.model';
import Empty from './Empty';
import * as S from './TestSpecs.styled';

interface IProps {
  assertionResults?: AssertionResults;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry, name: string): void;
  onOpen(selector: string): void;
  onRevert(originalSelector: string): void;
}

const TestSpecs = ({assertionResults, onDelete, onEdit, onOpen, onRevert}: IProps) => {
  if (!assertionResults?.resultList?.length) {
    return <Empty />;
  }

  return (
    <S.Container data-cy="test-specs-container">
      {assertionResults?.resultList?.map(specResult =>
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
