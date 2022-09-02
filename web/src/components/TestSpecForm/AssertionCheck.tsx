import {Form, FormInstance, Input, Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {capitalize} from 'lodash';
import {useMemo} from 'react';
import {durationRegExp} from '../../constants/Common.constants';
import {CompareOperator} from '../../constants/Operator.constants';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import OperatorService from '../../services/Operator.service';
import {TAssertion} from '../../types/Assertion.types';
import {TSpanFlatAttribute} from '../../types/Span.types';
import {DurationFields} from './DurationFields';
import {AttributeField} from './Fields/AttributeField';
import {OtelReference} from './hooks/useGetOTELSemanticConventionAttributesInfo';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';

const operatorList = Object.values(CompareOperator).map(value => ({
  value: OperatorService.getOperatorSymbol(value),
  label: capitalize(OperatorService.getOperatorName(value)),
}));

interface IProps {
  remove(name: number): void;
  reference: OtelReference;
  form: FormInstance<IValues>;
  field: Pick<FormListFieldData, never>;
  name: number;
  attributeList: TSpanFlatAttribute[];
  index: number;
  assertions: TAssertion[];
}

export const AssertionCheck = ({attributeList, field, index, name, assertions, form, remove, reference}: IProps) => {
  const assertion = assertions?.[index];
  const match = useMemo(() => assertion?.expected?.match(durationRegExp), [assertion?.expected]);

  return (
    <>
      <AttributeField field={field} name={name} attributeList={attributeList} reference={reference} />
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

      <div>
        {index !== 0 && (
          <S.DeleteCheckIcon
            onClick={() => {
              CreateAssertionModalAnalyticsService.onRemoveCheck();
              remove(name);
            }}
          />
        )}
      </div>
    </>
  );
};
