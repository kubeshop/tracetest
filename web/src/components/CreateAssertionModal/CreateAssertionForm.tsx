import React, {useCallback, useEffect, useMemo} from 'react';
import {isEmpty} from 'lodash';
import {QuestionCircleOutlined, PlusOutlined} from '@ant-design/icons';
import styled from 'styled-components';
import {Button, Input, Select as AntSelect, AutoComplete, Typography, Tooltip, Form, Space, FormInstance} from 'antd';
import jemsPath from 'jmespath';

import {COMPARE_OPERATOR, ISpan, ItemSelector, ITrace, LOCATION_NAME, SpanSelector} from 'types';
import {useCreateAssertionMutation} from 'services/TestService';
import {SELECTOR_DEFAULT_ATTRIBUTES} from 'lib/SelectorDefaultAttributes';
import {filterBySpanId} from 'utils';
import {CreateAssertionSelectorInput} from './CreateAssertionSelectorInput';
import {getSpanSignature} from '../../services/SpanService';
import {getSpanAttributeValueType} from '../../services/SpanAttributeService';

interface AssertionSpan {
  key: string;
  compareOp: keyof typeof COMPARE_OPERATOR;
  value: string;
}

const Select = styled(AntSelect)`
  min-width: 88px;
  > .ant-select-selector {
    min-height: 100%;
  }
`;

export type TValues = {
  assertionList: AssertionSpan[];
  selectorList: ItemSelector[];
};

const itemSelectorKeys = SELECTOR_DEFAULT_ATTRIBUTES.map(el => el.attributes).flat();

interface TCreateAssertionFormProps {
  onCreate(): void;
  onForm(form: FormInstance): void;
  onSelectorList(selectorList: ItemSelector[]): void;
  span: ISpan;
  testId: string;
  trace: ITrace;
}

const CreateAssertionForm: React.FC<TCreateAssertionFormProps> = ({
  testId,
  span,
  trace,
  onForm,
  onCreate,
  onSelectorList,
}) => {
  const [createAssertion] = useCreateAssertionMutation();
  const attrs = jemsPath.search(trace, filterBySpanId(span.spanId));
  const [form] = Form.useForm<TValues>();
  const spanSignature = useMemo(() => getSpanSignature(span.spanId, trace), [span.spanId, trace]);

  useEffect(() => {
    onForm(form);
  }, [onForm, form]);

  useEffect(() => {
    onSelectorList(spanSignature);
  }, [onSelectorList, spanSignature]);

  const spanTagsMap = attrs?.reduce((acc: {[x: string]: any}, item: {key: string}) => {
    if (itemSelectorKeys.indexOf(item.key) !== -1) {
      return acc;
    }
    const keyPrefix = item.key.split('.').shift() || item.key;
    if (!keyPrefix) {
      return acc;
    }
    const keys = acc[keyPrefix] || [];
    keys.push(item);
    acc[keyPrefix] = keys;
    return acc;
  }, {});

  const renderTitle = (title: any, index: number) => (
    <span key={`KEY_${title}_${index}`}>{`${title} (${spanTagsMap[title][0].type})`}</span>
  );

  const renderItem = (attr: any) => ({
    value: attr.key,
    label: (
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
        }}
      >
        {attr.key}
      </div>
    ),
  });

  const spanAssertionOptions = Object.keys(spanTagsMap).map((tagKey: any, index) => {
    return {
      label: renderTitle(tagKey, index),
      options: spanTagsMap[tagKey].map((el: any) => renderItem(el)),
    };
  });

  const handleCreateAssertion = useCallback(
    async ({assertionList, selectorList}: TValues) => {
      const spanAssertions = assertionList
        .map(k => ({value: k.value, attr: attrs?.find((el: any) => el.key === k?.key), compareOp: k.compareOp}))
        .filter(i => i.attr)
        .map(({attr, value, compareOp}) => {
          const spanAttribute = span.attributes.find(({key}) => key.includes(attr.key));

          return {
            locationName: attr.type.includes('span')
              ? LOCATION_NAME.SPAN_ATTRIBUTES
              : LOCATION_NAME.RESOURCE_ATTRIBUTES,
            propertyName: attr.key as string,
            comparisonValue: value,
            operator: compareOp,
            valueType: getSpanAttributeValueType(spanAttribute!),
          };
        })
        .filter((el): el is SpanSelector => Boolean(el));

      await createAssertion({testId, selectors: selectorList, spanAssertions}).unwrap();
      onCreate();
    },
    [attrs, createAssertion, onCreate, span.attributes, testId]
  );

  return (
    <Form
      name="newTest"
      form={form}
      initialValues={{
        remember: true,
        assertionList: [
          {
            key: '',
            compare: COMPARE_OPERATOR.EQUALS,
            value: '',
          },
        ],
        selectorList: spanSignature,
      }}
      onFinish={handleCreateAssertion}
      autoComplete="off"
      layout="vertical"
      onFieldsChange={changedFields => {
        onSelectorList(form.getFieldValue('selectorList') || []);
        const [field] = changedFields;

        const [fieldName = '', entry = 0, keyName = ''] = field.name as Array<string | number>;

        if (fieldName === 'assertionList' && keyName === 'key' && field.value) {
          const assertionList: AssertionSpan[] = form.getFieldValue('assertionList') || [];

          form.setFieldsValue({
            assertionList: assertionList.map((assertion, index) => {
              if (index === entry) {
                const value = attrs?.find((el: any) => el.key === assertionList[index].key)?.value;
                const isValid = typeof value === 'number' || !isEmpty(value);

                return {...assertion, value: isValid ? String(value) : ''};
              }

              return assertion;
            }),
          });
        }
      }}
    >
      <div style={{marginBottom: 8}}>
        <Typography.Text style={{marginRight: 8}}>Selectors</Typography.Text>
        <Tooltip title="Pick the attributes that filter the list of spans selected by this selector" placement="right">
          <QuestionCircleOutlined style={{color: '#8C8C8C'}} />
        </Tooltip>
      </div>
      <Form.Item name="selectorList" rules={[{required: true, message: 'At least one selector is required'}]}>
        <CreateAssertionSelectorInput spanSignature={spanSignature} />
      </Form.Item>
      <div style={{marginTop: 24, marginBottom: 8}}>
        <Typography.Text style={{marginRight: 8}}>Span Assertions</Typography.Text>
        <Tooltip title="Define the checks to run against each span selected by the list of selectors" placement="right">
          <QuestionCircleOutlined style={{color: '#8C8C8C'}} />
        </Tooltip>
      </div>
      <Form.List name="assertionList">
        {(fields, {add}) => {
          return (
            <>
              {fields.map(field => (
                <Space key={field.key} style={{display: 'flex', alignItems: 'stretch', gap: '4px', marginBottom: 16}}>
                  <Form.Item
                    {...field}
                    name={[field.name, 'key']}
                    style={{margin: 0}}
                    rules={[{required: true, message: 'Attribute is required'}]}
                  >
                    <AutoComplete
                      style={{width: 250, margin: 0}}
                      options={spanAssertionOptions}
                      filterOption={(inputValue, option) => {
                        return option?.label.props.children.includes(inputValue);
                      }}
                    >
                      <Input.Search size="large" placeholder="span key" />
                    </AutoComplete>
                  </Form.Item>
                  <Form.Item
                    {...field}
                    initialValue={COMPARE_OPERATOR.EQUALS}
                    style={{margin: 0}}
                    name={[field.name, 'compareOp']}
                    rules={[{required: true, message: 'Operator is required'}]}
                  >
                    <Select style={{margin: 0}}>
                      <Select.Option value={COMPARE_OPERATOR.EQUALS}>eq</Select.Option>
                      <Select.Option value={COMPARE_OPERATOR.NOTEQUALS}>ne</Select.Option>
                      <Select.Option value={COMPARE_OPERATOR.GREATERTHAN}>gt</Select.Option>
                      <Select.Option value={COMPARE_OPERATOR.LESSTHAN}>lt</Select.Option>
                      <Select.Option value={COMPARE_OPERATOR.GREATOREQUALS}>ge</Select.Option>
                      <Select.Option value={COMPARE_OPERATOR.LESSOREQUAL}>le</Select.Option>
                    </Select>
                  </Form.Item>
                  <Form.Item
                    {...field}
                    name={[field.name, 'value']}
                    style={{margin: 0}}
                    rules={[{required: true, message: 'Value is required'}]}
                  >
                    <Input placeholder="value" />
                  </Form.Item>
                </Space>
              ))}

              <Button type="link" icon={<PlusOutlined />} style={{padding: 0}} onClick={add}>
                Add Item
              </Button>
            </>
          );
        }}
      </Form.List>
    </Form>
  );
};

export default CreateAssertionForm;
