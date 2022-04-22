import React, {useCallback, useEffect, useMemo} from 'react';
import {isEmpty} from 'lodash';
import {QuestionCircleOutlined, PlusOutlined, MinusCircleOutlined} from '@ant-design/icons';
import {Button, Input, AutoComplete, Typography, Tooltip, Form, Space, FormInstance} from 'antd';
import jemsPath from 'jmespath';

import {Assertion, COMPARE_OPERATOR, ISpan, ItemSelector, ITrace, LOCATION_NAME, SpanSelector} from 'types';
import {useCreateAssertionMutation, useUpdateAssertionMutation} from 'services/TestService';
import {SELECTOR_DEFAULT_ATTRIBUTES} from 'lib/SelectorDefaultAttributes';
import {filterBySpanId} from 'utils';
import {CreateAssertionSelectorInput} from './CreateAssertionSelectorInput';
import {getSpanSignature} from '../../services/SpanService';
import {getSpanAttributeValueType} from '../../services/SpanAttributeService';
import * as S from './CreateAssertionModal.styled';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTourService';
import {Steps} from '../GuidedTour/assertionStepList';
import useGuidedTour from '../GuidedTour/useGuidedTour';

interface AssertionSpan {
  key: string;
  compareOp: keyof typeof COMPARE_OPERATOR;
  value: string;
}

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
  assertion?: Assertion;
}

const CreateAssertionForm: React.FC<TCreateAssertionFormProps> = ({
  testId,
  span,
  trace,
  assertion,
  onForm,
  onCreate,
  onSelectorList,
}) => {
  const [createAssertion] = useCreateAssertionMutation();
  const [updateAssertion] = useUpdateAssertionMutation();
  const attrs = jemsPath.search(trace, filterBySpanId(span.spanId));
  const [form] = Form.useForm<TValues>();
  const defaultSelectorList = useMemo(() => getSpanSignature(span.spanId, trace), [span.spanId, trace]);

  useGuidedTour(GuidedTours.Assertion);

  const defaultAssertionList = useMemo<AssertionSpan[]>(() => {
    if (assertion) {
      return assertion.spanAssertions.map(({propertyName, operator, comparisonValue}) => ({
        key: propertyName,
        compareOp: operator,
        value: comparisonValue,
      }));
    }

    return [
      {
        key: '',
        compareOp: COMPARE_OPERATOR.EQUALS,
        value: '',
      },
    ];
  }, [assertion]);

  useEffect(() => {
    onForm(form);
  }, [onForm, form]);

  useEffect(() => {
    onSelectorList(assertion ? assertion.selectors : defaultSelectorList);
  }, [onSelectorList, defaultSelectorList]);

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

      const newData = {selectors: selectorList, spanAssertions};

      if (assertion) {
        await updateAssertion({testId, assertionId: assertion.assertionId, assertion: newData});
      } else {
        await createAssertion({testId, assertion: newData});
      }
      onCreate();
    },
    [assertion, attrs, createAssertion, onCreate, span.attributes, testId, updateAssertion]
  );

  return (
    <Form
      name="newTest"
      form={form}
      initialValues={{
        remember: true,
        assertionList: defaultAssertionList,
        selectorList: assertion ? assertion.selectors : defaultSelectorList,
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
            assertionList: assertionList.map((assertionEntry, index) => {
              if (index === entry) {
                const value = attrs?.find((el: any) => el.key === assertionList[index].key)?.value;
                const isValid = typeof value === 'number' || !isEmpty(value);

                return {...assertionEntry, value: isValid ? String(value) : ''};
              }

              return assertionEntry;
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
      <Form.Item
        name="selectorList"
        rules={[{required: true, message: 'At least one selector is required'}]}
        data-tour={GuidedTourService.getStep(GuidedTours.Assertion, Steps.Selectors)}
      >
        <CreateAssertionSelectorInput spanSignature={defaultSelectorList} />
      </Form.Item>
      <div style={{marginTop: 24, marginBottom: 8}}>
        <Typography.Text style={{marginRight: 8}}>Span Assertions</Typography.Text>
        <Tooltip title="Define the checks to run against each span selected by the list of selectors" placement="right">
          <QuestionCircleOutlined style={{color: '#8C8C8C'}} />
        </Tooltip>
      </div>
      <div data-tour={GuidedTourService.getStep(GuidedTours.Assertion, Steps.Checks)}>
        <Form.List name="assertionList">
          {(fields, {add, remove}) => {
            return (
              <>
                {fields.map(({key, name, ...field}, index) => (
                  <Space key={key} style={{display: 'flex', alignItems: 'center', gap: '4px', marginBottom: 16}}>
                    <Form.Item
                      {...field}
                      name={[name, 'key']}
                      style={{margin: 0}}
                      rules={[{required: true, message: 'Attribute is required'}]}
                    >
                      <AutoComplete
                        style={{margin: 0}}
                        options={spanAssertionOptions}
                        filterOption={(inputValue, option) => {
                          return option?.label.props.children.includes(inputValue);
                        }}
                      >
                        <Input.Search size="large" placeholder="span key" />
                      </AutoComplete>
                    </Form.Item>
                    <S.FullHeightFormItem
                      {...field}
                      initialValue={COMPARE_OPERATOR.EQUALS}
                      style={{margin: 0}}
                      name={[name, 'compareOp']}
                      rules={[{required: true, message: 'Operator is required'}]}
                    >
                      <S.Select style={{margin: 0}}>
                        <S.Select.Option value={COMPARE_OPERATOR.EQUALS}>eq</S.Select.Option>
                        <S.Select.Option value={COMPARE_OPERATOR.NOTEQUALS}>ne</S.Select.Option>
                        <S.Select.Option value={COMPARE_OPERATOR.GREATERTHAN}>gt</S.Select.Option>
                        <S.Select.Option value={COMPARE_OPERATOR.LESSTHAN}>lt</S.Select.Option>
                        <S.Select.Option value={COMPARE_OPERATOR.GREATOREQUALS}>ge</S.Select.Option>
                        <S.Select.Option value={COMPARE_OPERATOR.LESSOREQUAL}>le</S.Select.Option>
                      </S.Select>
                    </S.FullHeightFormItem>
                    <S.FullHeightFormItem
                      {...field}
                      name={[name, 'value']}
                      style={{margin: 0}}
                      rules={[{required: true, message: 'Value is required'}]}
                    >
                      <Input placeholder="value" />
                    </S.FullHeightFormItem>
                    <MinusCircleOutlined
                      color="error"
                      style={{cursor: 'pointer', color: 'rgb(140, 140, 140)'}}
                      onClick={() => index > 0 && remove(name)}
                    />
                  </Space>
                ))}

                <Button type="link" icon={<PlusOutlined />} style={{padding: 0}} onClick={add}>
                  Add Item
                </Button>
              </>
            );
          }}
        </Form.List>
      </div>
    </Form>
  );
};

export default CreateAssertionForm;
