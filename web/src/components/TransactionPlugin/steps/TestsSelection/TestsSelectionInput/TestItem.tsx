import {SortableElement} from 'react-sortable-hoc';
import {TTest} from 'types/Test.types';
import * as S from './TestsSelectionInput.styled';

interface IProps {
  value: TTest;
  onDelete(testId: string): void;
}

const TestItem = ({value: test, onDelete}: IProps) => {
  return (
    <S.TestItemContainer>
      <S.DragHandle />
      <span>{test.name}</span>
      <S.DeleteIcon onClick={() => onDelete(test.id)} />
    </S.TestItemContainer>
  );
};

export default SortableElement(TestItem);
