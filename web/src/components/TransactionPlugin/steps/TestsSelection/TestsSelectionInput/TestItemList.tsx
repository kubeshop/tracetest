import {Col, Row} from 'antd';
import {TTest} from 'types/Test.types';
import TestItem from './TestItem';
import * as S from './TestsSelectionInput.styled';

interface IProps {
  items: TTest[];
  onDelete(testId: string): void;
}

const TestItemList = ({items, onDelete}: IProps) => {
  return (
    <Row gutter={12}>
      <Col span={18}>
        <S.ItemListContainer>
          {items.map((test, index) => (
            // eslint-disable-next-line react/no-array-index-key
            <TestItem key={`${test.id}-${index}`} test={test} onDelete={onDelete} />
          ))}
        </S.ItemListContainer>
      </Col>
    </Row>
  );
};

export default TestItemList;
