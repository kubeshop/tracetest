import {Form, FormInstance, Input, Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {capitalize} from 'lodash';
import {useMemo} from 'react';

import AttributesTags from '../../constants/AttributesTags.json';
import {durationRegExp} from '../../constants/Common.constants';
import {CompareOperator} from '../../constants/Operator.constants';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import OperatorService from '../../services/Operator.service';
import {TAssertion} from '../../types/Assertion.types';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';
import {DurationFields} from './DurationFields';
import {OtelReference} from './hooks/useGetOTELSemanticConventionAttributesInfo';

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
  attributeOptionList: JSX.Element[];
  index: number;
  assertionList: TAssertion[];
}

export const AssertionCheck = ({
  attributeOptionList,
  field,
  index,
  name,
  assertionList,
  form,
  remove,
  reference,
}: IProps) => {
  const assertion = assertionList?.[index];
  const match = useMemo(() => assertion?.expected?.match(durationRegExp), [assertion?.expected]);
  return (
    <>
      <Form.Item
        {...field}
        name={[name, 'attribute']}
        style={{margin: 0}}
        rules={[{required: true, message: 'Attribute is required'}]}
        data-cy="assertion-check-attribute"
        id="assertion-check-attribute"
      >
        <S.Select
          style={{margin: 0}}
          placeholder="Attribute"
          showSearch
          filterOption={(input, option) => {
            const key = option?.key || '';
            const attributesTags: Record<string, {tags: string[]; description: string}> = AttributesTags;
            const ref = reference[key] || attributesTags[key] || {description: '', tags: []};
            const availableTagsMatchInput = Boolean(
              ref.tags.find(tag => tag.toString().toLowerCase().includes(input.toLowerCase()))
            );
            const currentOptionMatchInput = option?.key.includes(input);
            const currentDescriptionMatchInput = ref?.description.includes(input);
            return availableTagsMatchInput || currentOptionMatchInput || currentDescriptionMatchInput;
          }}
        >
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
