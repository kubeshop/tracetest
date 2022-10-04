import {useSortable} from '@dnd-kit/sortable';
import {CSS} from '@dnd-kit/utilities';
import {TTest} from 'types/Test.types';
import * as S from './TestsSelectionInput.styled';

interface IProps {
  test: TTest;
  onDelete(testId: string): void;
}

const TestItem = ({test, onDelete}: IProps) => {
  const {attributes, listeners, setNodeRef, transform, transition} = useSortable({id: test.id});

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <S.TestItemContainer ref={setNodeRef} style={style} {...attributes}>
      <S.DragHandle {...listeners} />
      <span>{test.name}</span>
      <S.DeleteIcon onClick={() => onDelete(test.id)} />
    </S.TestItemContainer>
  );
};

export default TestItem;
