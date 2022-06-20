import {Typography, Form, Button} from 'antd';
import {FieldData} from 'antd/node_modules/rc-field-form/es/interface';
import {isEmpty} from 'lodash';
import React, {useCallback, useEffect} from 'react';

import {CompareOperator} from 'constants/Operator.constants';
import {useGetSelectedSpansQuery} from 'redux/apis/TraceTest.api';
import {useAppSelector} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import OperatorService from 'services/Operator.service';
import SelectorService from 'services/Selector.service';
import {TAssertion, TPseudoSelector, TSpanSelector} from 'types/Assertion.types';
import {TooltipQuestion} from '../TooltipQuestion/TooltipQuestion';
import * as S from './AssertionForm.styled';
import AssertionFormCheckList from './AssertionFormCheckList';
import AssertionFormPseudoSelectorInput from './AssertionFormPseudoSelectorInput';
import AssertionFormSelectorInput from './AssertionFormSelectorInput';
import {useSpan} from '../../providers/Span/Span.provider';

const {onChecksChange, onSelectorChange} = CreateAssertionModalAnalyticsService;

export interface IValues {
  assertionList?: TAssertion[];
  selectorList: TSpanSelector[];
  pseudoSelector?: TPseudoSelector;
}

interface TAssertionFormProps {
  defaultValues?: IValues;
  onSubmit(values: IValues): void;
  testId: string;
  runId: string;
  isEditing?: boolean;
  onCancel(): void;
}

const AssertionForm: React.FC<TAssertionFormProps> = ({
  defaultValues: {
    assertionList = [
      {
        comparator: OperatorService.getOperatorSymbol(CompareOperator.EQUALS),
      },
    ],
    selectorList = [],
    pseudoSelector,
  } = {},
  onSubmit,
  onCancel,
  isEditing = false,
  testId,
  runId,
}) => {
  const {onSetAffectedSpans, onClearAffectedSpans} = useSpan();
  const [form] = Form.useForm<IValues>();

  const currentSelectorList = Form.useWatch('selectorList', form) || [];
  const currentAssertionList = Form.useWatch('assertionList', form) || [];
  const currentPseudoSelector = Form.useWatch('pseudoSelector', form) || undefined;

  const {data: spanIdList = []} = useGetSelectedSpansQuery({
    query: SelectorService.getSelectorString(currentSelectorList, currentPseudoSelector),
    testId,
    runId,
  });

  useEffect(() => {
    onSetAffectedSpans(spanIdList);

    return () => {
      onClearAffectedSpans();
    };
  }, [onClearAffectedSpans, onSetAffectedSpans, spanIdList]);

  const attributeList = useAppSelector(state =>
    AssertionSelectors.selectAttributeList(state, testId, runId, spanIdList)
  );
  const selectorAttributeList = useAppSelector(state =>
    AssertionSelectors.selectSelectorAttributeList(state, testId, runId, spanIdList, currentSelectorList)
  );
  const definitionSelectorList = useAppSelector(state => TestDefinitionSelectors.selectDefinitionSelectorList(state));

  const onFieldsChange = useCallback(
    (changedFields: FieldData[]) => {
      const [field] = changedFields;

      const [fieldName = '', entry = 0, keyName = ''] = field.name as Array<string | number>;

      if (fieldName === 'selectorList') onSelectorChange();
      if (fieldName === 'assertionList') onChecksChange();

      if (fieldName === 'assertionList' && keyName === 'attribute' && field.value) {
        const list: TAssertion[] = form.getFieldValue('assertionList') || [];

        form.setFieldsValue({
          assertionList: list.map((assertionEntry, index) => {
            if (index === entry) {
              const {value = ''} = attributeList?.find((el: any) => el.key === list[index].attribute) || {};
              const isValid = typeof value === 'number' || !isEmpty(value);

              return {...assertionEntry, expected: isValid ? String(value) : ''};
            }

            return assertionEntry;
          }),
        });
      }
    },
    [attributeList, form]
  );

  return (
    <S.AssertionForm>
      <S.AssertionFormHeader>
        <S.AssertionFormTitle strong>{isEditing ? 'Edit Assertion' : 'Add New Assertion'}</S.AssertionFormTitle>
        <Typography.Text type="secondary">Affects {spanIdList.length} span(s)</Typography.Text>
      </S.AssertionFormHeader>
      <Form
        name="assertion-form"
        form={form}
        initialValues={{
          remember: true,
          assertionList,
          selectorList,
          pseudoSelector,
        }}
        onFinish={onSubmit}
        autoComplete="off"
        layout="vertical"
        data-cy="assertion-form"
        onFieldsChange={onFieldsChange}
      >
        <div style={{marginBottom: 8}}>
          <Typography.Text>Filter to limit the span(s) included in this assertion</Typography.Text>
          <TooltipQuestion
            title={`
            You can decided which spans will be tested by this assertion by altering the filter. 
            Use the dropdown to the right to select the first matching span, last, n-th, or all.  
            `}
          />
        </div>
        <S.SelectorInputContainer>
          <Form.Item
            name="selectorList"
            rules={[
              {required: true, message: 'At least one selector is required'},
              {
                validator: (_, value: TSpanSelector[]) =>
                  SelectorService.validateSelector(
                    definitionSelectorList,
                    isEditing,
                    selectorList,
                    value,
                    currentPseudoSelector,
                    pseudoSelector
                  ),
              },
            ]}
          >
            <AssertionFormSelectorInput attributeList={selectorAttributeList} />
          </Form.Item>
          <Form.Item name="pseudoSelector">
            <AssertionFormPseudoSelectorInput />
          </Form.Item>
        </S.SelectorInputContainer>
        <div style={{marginBottom: 8}}>
          <Typography.Text>Define the checks to run against each span selected</Typography.Text>
          <TooltipQuestion
            title={`
            Add one of more checks to be run against the span(s) that match your filter.  
            For example, create one assertion to check all http spans to make sure they return status code 200... 
            all in one assertion.
            `}
          />
        </div>
        <div>
          <Form.List name="assertionList">
            {(fields, {add, remove}) => (
              <AssertionFormCheckList
                assertionList={currentAssertionList}
                fields={fields}
                add={add}
                remove={remove}
                attributeList={attributeList}
              />
            )}
          </Form.List>
        </div>
        <S.AssertionFromActions>
          <Button onClick={onCancel}>Cancel</Button>
          <Button type="primary" onClick={form.submit} data-cy="assertion-form-submit-button">
            {isEditing ? 'Save' : 'Add'}
          </Button>
        </S.AssertionFromActions>
      </Form>
    </S.AssertionForm>
  );
};

export default AssertionForm;
