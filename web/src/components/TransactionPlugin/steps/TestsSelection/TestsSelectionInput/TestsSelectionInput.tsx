import {useState, useCallback} from 'react';
import {arrayMoveImmutable} from 'array-move';
import {Col, Row, Select} from 'antd';
import {noop} from 'lodash';
import {TTest} from 'types/Test.types';
import TestItemList from './TestItemList';

interface IProps {
  onChange?(tests: string[]): void;
  value?: string[];
  testList: TTest[];
}

const TestsSelectionInput = ({value = [], onChange = noop, testList}: IProps) => {
  const [selectedTestList, setSelectedTestList] = useState<TTest[]>([]);
  const onSelectedTest = useCallback(
    (testId: string) => {
      onChange([...value, testId]);
      setSelectedTestList([...selectedTestList, testList.find(test => test.id === testId)!]);
    },
    [onChange, selectedTestList, testList, value]
  );

  const onSortEnd = useCallback(
    ({oldIndex, newIndex}) => {
      const updatedList = arrayMoveImmutable(selectedTestList, oldIndex, newIndex);
      setSelectedTestList(updatedList);

      onChange(updatedList.map(test => test.id));
    },
    [onChange, selectedTestList]
  );

  const onDelete = useCallback(
    (testId: string) => {
      onChange(value.filter(id => id !== testId));
      setSelectedTestList(selectedTestList.filter(test => test.id !== testId));
    },
    [onChange, selectedTestList, value]
  );

  return (
    <>
      <TestItemList
        items={selectedTestList}
        onSortEnd={onSortEnd}
        helperClass="draggable-item"
        onDelete={onDelete}
        useDragHandle
      />
      <Row gutter={12}>
        <Col span={18}>
          <span>
            <Select
              placeholder="Add a test"
              onChange={onSelectedTest}
              value={null}
              showSearch
              filterOption={(input, option) =>
                (option!.children as unknown as string).toLowerCase().includes(input.toLowerCase())
              }
            >
              {testList.map(({id, name}) => (
                <Select.Option value={id} key={id}>
                  {name}
                </Select.Option>
              ))}
            </Select>
          </span>
        </Col>
      </Row>
    </>
  );
};

export default TestsSelectionInput;
