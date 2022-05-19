import {PlusOutlined} from '@ant-design/icons';
import {Input, Select, Form} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {capitalize, uniqBy} from 'lodash';
import {useMemo} from 'react';
import {CompareOperator} from '../../constants/Operator.constants';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import OperatorService from '../../services/Operator.service';
import {IAssertionSpan} from '../../types/Assertion.types';
import {ISpanFlatAttribute} from '../../types/Span.types';
import * as S from './AssertionForm.styled';

interface IProps {
  fields: FormListFieldData[];
  assertionList: IAssertionSpan[];
  add(): void;
  remove(name: number): void;
  attributeList: ISpanFlatAttribute[];
}

const operatorList = Object.values(CompareOperator).map(value => ({
  value,
  label: capitalize(OperatorService.getOperatorName(value)),
}));

const getIsValid = ({key, compareOp, value}: IAssertionSpan): boolean => Boolean(key && compareOp && value);

const AssertionFormCheckList: React.FC<IProps> = ({fields, add, remove, attributeList, assertionList}) => {
  const attributeOptionList = useMemo(
    () =>
      uniqBy(attributeList, 'key').map(({key}) => (
        <Select.Option key={key} value={key}>
          {key}
        </Select.Option>
      )),
    [attributeList]
  );

  return (
    <>
      {fields.map(({key, name, ...field}, index) => (
        <S.Check key={key}>
          <Form.Item
            {...field}
            name={[name, 'key']}
            style={{margin: 0}}
            rules={[{required: true, message: 'Attribute is required'}]}
            data-cy="assertion-check-attribute"
            id="assertion-check-attribute"
          >
            <S.Select style={{margin: 0}} placeholder="Attribute" showSearch>
              {attributeOptionList}
            </S.Select>
          </Form.Item>
          <Form.Item
            {...field}
            style={{margin: 0}}
            name={[name, 'compareOp']}
            rules={[{required: true, message: 'Operator is required'}]}
            data-cy="assertion-check-operator"
          >
            <S.Select style={{margin: 0}} placeholder="Assertion Type">
              {operatorList.map(({value, label}) => (
                <Select.Option key={value} value={value}>
                  {label}
                </Select.Option>
              ))}
            </S.Select>
          </Form.Item>
          <Form.Item
            {...field}
            name={[name, 'value']}
            style={{margin: 0}}
            rules={[{required: true, message: 'Value is required'}]}
            data-cy="assertion-check-value"
          >
            <Input placeholder="Expected Value" />
          </Form.Item>
          {index === fields.length - 1 ? (
            <S.AddCheckButton
              icon={<PlusOutlined />}
              onClick={() => {
                CreateAssertionModalAnalyticsService.onAddCheck();
                add();
              }}
              disabled={!assertionList[name] || !getIsValid(assertionList[name])}
              data-cy="add-assertion-form-add-check"
            >
              Add Check
            </S.AddCheckButton>
          ) : (
            <S.DeleteCheckIcon
              onClick={() => {
                CreateAssertionModalAnalyticsService.onRemoveCheck();
                remove(name);
              }}
            />
          )}
        </S.Check>
      ))}
    </>
  );
};

export default AssertionFormCheckList;
