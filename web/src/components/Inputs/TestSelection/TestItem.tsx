import {useSortable} from '@dnd-kit/sortable';
import {CSS} from '@dnd-kit/utilities';
import Test from 'models/Test.model';
import * as S from './TestsSelection.styled';

interface IProps {
  test: Test;
  sortableId: string;
  onDelete(sortableId: string): void;
}

const TestItem = ({test, onDelete, sortableId}: IProps) => {
  const {attributes, listeners, setNodeRef, transform, transition} = useSortable({id: sortableId});

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <S.TestItemContainer ref={setNodeRef} style={style} {...attributes}>
      <S.DragHandle {...listeners} />
      <S.TestNameContainer>
        <S.TestLink to={`/test/${test.id}`} target="_blank">
          <span>{test.name}</span>
        </S.TestLink>
      </S.TestNameContainer>
      <S.DeleteIcon onClick={() => onDelete(sortableId)} />
    </S.TestItemContainer>
  );
};

export default TestItem;
