import {Col, Row} from 'antd';
import TestItem from './TestItem';
import {ISortableTest} from './TestsSelectionInput';
import * as S from './TestsSelectionInput.styled';

interface IProps {
  items: ISortableTest[];
  onDelete(sortableId: string): void;
}

const TestItemList = ({items, onDelete}: IProps) => {
  return (
    <Row gutter={12}>
      <Col span={18}>
        <S.ItemListContainer>
          {items.map(({id, test}) => (
            <TestItem key={id} test={test} sortableId={id} onDelete={onDelete} />
          ))}
        </S.ItemListContainer>
      </Col>
    </Row>
  );
};

export default TestItemList;
