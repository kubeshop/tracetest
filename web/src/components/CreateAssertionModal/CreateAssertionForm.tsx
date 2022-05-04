import React, {useCallback, useEffect, useMemo} from 'react';
import {isEmpty} from 'lodash';
import {QuestionCircleOutlined, PlusOutlined, MinusCircleOutlined} from '@ant-design/icons';
import {Button, Input, AutoComplete, Typography, Tooltip, Form, Space, FormInstance} from 'antd';

import {useCreateAssertionMutation, useUpdateAssertionMutation} from 'redux/apis/Test.api';
import {SELECTOR_DEFAULT_ATTRIBUTES} from 'constants/SemanticGroupNames.constants';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import useGuidedTour from 'hooks/useGuidedTour';
import {CreateAssertionSelectorInput} from './CreateAssertionSelectorInput';
import * as S from './CreateAssertionModal.styled';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import {IAssertion, IItemSelector, ISpanSelector} from '../../types/Assertion.types';
import {CompareOperator} from '../../constants/Operator.constants';
import {ISpan} from '../../types/Span.types';
import {LOCATION_NAME} from '../../constants/Span.constants';
import {Steps} from '../GuidedTour/assertionStepList';

const {
  onAddCheck,
  onRemoveCheck,
  onChecksChange,
  onCreateAssertionFormSubmit,
  onEditAssertionFormSubmit,
  onSelectorChange,
} = CreateAssertionModalAnalyticsService;

interface IAssertionSpan {
  key: string;
  compareOp: CompareOperator;
  value: string;
}

export type TValues = {
  assertionList: IAssertionSpan[];
  selectorList: IItemSelector[];
};

const itemSelectorKeys = SELECTOR_DEFAULT_ATTRIBUTES.map(el => el.attributes).flat();

interface TCreateAssertionFormProps {
  onCreate(): void;
  onForm(form: FormInstance): void;
  onSelectorList(selectorList: IItemSelector[]): void;
  span: ISpan;
  testId: string;
  assertion?: IAssertion;
}

const CreateAssertionForm: React.FC<TCreateAssertionFormProps> = ({
  testId,
  span,
  assertion,
  onForm,
  onCreate,
  onSelectorList,
}) => {
  const [createAssertion] = useCreateAssertionMutation();
  const [updateAssertion] = useUpdateAssertionMutation();
  const [form] = Form.useForm<TValues>();
  const {attributeList, signature: defaultSelectorList} = span;

  useGuidedTour(GuidedTours.Assertion);

  const defaultAssertionList = useMemo<IAssertionSpan[]>(() => {
    if (assertion) {
      return assertion.spanAssertions?.map(({propertyName, operator, comparisonValue}) => ({
        key: propertyName,
        compareOp: operator,
        value: comparisonValue,
      }));
    }

    return [
      {
        key: '',
        compareOp: CompareOperator.EQUALS,
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

  const spanTagsMap = attributeList?.reduce((acc: {[x: string]: any}, item: {key: string}) => {
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
        .map(({value, compareOp, key}) => {
          const {name, type} = span.attributes[key];

          return {
            locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
            propertyName: name,
            comparisonValue: value,
            operator: compareOp,
            valueType: type,
          };
        })
        .filter((el): el is ISpanSelector => Boolean(el));

      const newData = {selectors: selectorList, spanAssertions};

      if (assertion) {
        onEditAssertionFormSubmit(assertion.assertionId);
        await updateAssertion({testId, assertionId: assertion.assertionId, assertion: newData});
      } else {
        onCreateAssertionFormSubmit(testId);
        await createAssertion({testId, assertion: newData});
      }
      onCreate();
    },
    [assertion, createAssertion, onCreate, span.attributes, testId, updateAssertion]
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
        const selectorList = form.getFieldValue('selectorList') || [];
        onSelectorList(selectorList);
        const [field] = changedFields;

        const [fieldName = '', entry = 0, keyName = ''] = field.name as Array<string | number>;

        if (fieldName === 'selectorList') onSelectorChange(JSON.stringify(selectorList));
        if (fieldName === 'assertionList') onChecksChange(JSON.stringify(form.getFieldValue('assertionList') || []));

        if (fieldName === 'assertionList' && keyName === 'key' && field.value) {
          const assertionList: IAssertionSpan[] = form.getFieldValue('assertionList') || [];

          form.setFieldsValue({
            assertionList: assertionList.map((assertionEntry, index) => {
              if (index === entry) {
                const value = attributeList?.find((el: any) => el.key === assertionList[index].key)?.value;
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
                {fields.map(({key, name, ...field}) => (
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
                      initialValue={CompareOperator.EQUALS}
                      style={{margin: 0}}
                      name={[name, 'compareOp']}
                      rules={[{required: true, message: 'Operator is required'}]}
                    >
                      <S.Select style={{margin: 0}}>
                        <S.Select.Option value={CompareOperator.EQUALS}>eq</S.Select.Option>
                        <S.Select.Option value={CompareOperator.NOTEQUALS}>ne</S.Select.Option>
                        <S.Select.Option value={CompareOperator.GREATERTHAN}>gt</S.Select.Option>
                        <S.Select.Option value={CompareOperator.LESSTHAN}>lt</S.Select.Option>
                        <S.Select.Option value={CompareOperator.GREATOREQUALS}>ge</S.Select.Option>
                        <S.Select.Option value={CompareOperator.LESSOREQUAL}>le</S.Select.Option>
                        <S.Select.Option value={CompareOperator.CONTAINS}>contains</S.Select.Option>
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
                      onClick={() => {
                        onRemoveCheck();
                        remove(name);
                      }}
                    />
                  </Space>
                ))}

                <Button
                  type="link"
                  icon={<PlusOutlined />}
                  style={{padding: 0}}
                  onClick={() => {
                    onAddCheck();
                    add();
                  }}
                >
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
