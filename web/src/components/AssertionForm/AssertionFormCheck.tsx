import {PlusOutlined} from '@ant-design/icons';
import {Form, FormInstance, Input, Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {capitalize} from 'lodash';
import {useMemo} from 'react';
import {durationRegExp} from '../../constants/Common.constants';
import {CompareOperator} from '../../constants/Operator.constants';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import OperatorService from '../../services/Operator.service';
import {TAssertion} from '../../types/Assertion.types';
import {IValues} from './AssertionForm';
import * as S from './AssertionForm.styled';
import {DurationFields} from './DurationFields';

const operatorList = Object.values(CompareOperator).map(value => ({
  value: OperatorService.getOperatorSymbol(value),
  label: capitalize(OperatorService.getOperatorName(value)),
}));

const getIsValid = ({attribute, comparator, expected}: TAssertion): boolean =>
  Boolean(attribute && comparator && expected);

interface IProps {
  add(): void;
  remove(name: number): void;
  form: FormInstance<IValues>;
  field: Pick<FormListFieldData, never>;
  name: number;
  attributeOptionList: JSX.Element[];
  index: number;
  length: number;
  assertionList: TAssertion[];
}

export const AssertionFormCheck = ({
  attributeOptionList,
  field,
  index,
  length,
  name,
  assertionList,
  form,
  add,
  remove,
}: IProps) => {
  const assertion = assertionList?.[index];
  const match = useMemo(() => assertion?.expected?.match(durationRegExp), [assertion?.expected]);
  return (
    <S.Check>
      <Form.Item
        {...field}
        name={[name, 'attribute']}
        style={{margin: 0}}
        rules={[{required: true, message: 'Attribute is required'}]}
        data-cy="assertion-check-attribute"
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
        initialValue={operatorList[0].value}
      >
        <S.Select style={{margin: 0}} placeholder="Assertion Type">
          {operatorList.map(({value, label}) => (
            <Select.Option key={value} value={value}>
              {label}
            </Select.Option>
          ))}
        </S.Select>
      </Form.Item>
      {match ? (
        <DurationFields form={form} index={index} assertion={assertion} />
      ) : (
        <Form.Item
          {...field}
          name={[name, 'expected']}
          style={{margin: 0}}
          rules={[{required: true, message: 'Value is required'}]}
          data-cy="assertion-check-value"
        >
          <Input placeholder="Expected Value" />
        </Form.Item>
      )}
      <S.CheckActions>
        {index !== 0 && (
          <S.DeleteCheckIcon
            onClick={() => {
              CreateAssertionModalAnalyticsService.onRemoveCheck();
              remove(name);
            }}
          />
        )}
        {index === length - 1 && (
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
  );
};
