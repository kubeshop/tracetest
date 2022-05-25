import {PlusOutlined} from '@ant-design/icons';
import {Input, Select, Form} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {capitalize, uniqBy} from 'lodash';
import {useMemo} from 'react';
import {CompareOperator} from '../../constants/Operator.constants';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import OperatorService from '../../services/Operator.service';
import {TAssertion} from '../../types/Assertion.types';
import {TSpanFlatAttribute} from '../../types/Span.types';
import * as S from './AssertionForm.styled';

interface IProps {
  fields: FormListFieldData[];
  assertionList: TAssertion[];

  add(): void;

  remove(name: number): void;

  attributeList: TSpanFlatAttribute[];
}

const operatorList = Object.values(CompareOperator).map(value => ({
  value: OperatorService.getOperatorSymbol(value),
  label: capitalize(OperatorService.getOperatorName(value)),
}));

const getIsValid = ({attribute, comparator, expected}: TAssertion): boolean =>
  Boolean(attribute && comparator && expected);

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
            name={[name, 'attribute']}
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
            name={[name, 'comparator']}
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
            name={[name, 'expected']}
            style={{margin: 0}}
            rules={[{required: true, message: 'Value is required'}]}
            data-cy="assertion-check-value"
          >
            <Input placeholder="Expected Value" />
          </Form.Item>

          <S.CheckActions>
            <S.DeleteCheckIcon
              onClick={() => {
                CreateAssertionModalAnalyticsService.onRemoveCheck();
                remove(name);
              }}
            />
            {index === fields.length - 1 && (
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
            )}
          </S.CheckActions>
        </S.Check>
      ))}
    </>
  );
};

export default AssertionFormCheckList;
