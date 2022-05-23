import React, {useCallback} from 'react';
import {useSelector} from 'react-redux';
import {FieldData} from 'antd/node_modules/rc-field-form/es/interface';
import {isEmpty} from 'lodash';
import {Typography, Form, Button} from 'antd';

import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import useGuidedTour from 'hooks/useGuidedTour';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import {CompareOperator, PseudoSelector} from '../../constants/Operator.constants';
import {Steps} from '../GuidedTour/assertionStepList';
import * as S from './AssertionForm.styled';
import AssertionSelectors from '../../selectors/Assertion.selectors';
import AssertionFormSelectorInput from './AssertionFormSelectorInput';
import {TAssertion, TSpanSelector} from '../../types/Assertion.types';
import AssertionFormPseudoSelectorInput from './AssertionFormPseudoSelectorInput';
import AssertionFormCheckList from './AssertionFormCheckList';

const {onChecksChange, onSelectorChange} = CreateAssertionModalAnalyticsService;

export interface IValues {
  assertionList: TAssertion[];
  selectorList: TSpanSelector[];
  pseudoSelector?: {
    selector: PseudoSelector;
    number?: number;
  };
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
        key: '',
        compareOp: CompareOperator.EQUALS,
        value: '',
      },
    ],
    selectorList = [],
  } = {},
  onSubmit,
  onCancel,
  isEditing = false,
  testId,
  runId,
}) => {
  const [form] = Form.useForm<IValues>();

  useGuidedTour(GuidedTours.Assertion);

  const currentSelectorList = Form.useWatch('selectorList', form) || [];
  const currentAssertionList = Form.useWatch('assertionList', form) || [];
  const attributeList = useSelector(AssertionSelectors.selectAttributeList(testId, runId, currentSelectorList));

  const onFieldsChange = useCallback(
    (changedFields: FieldData[]) => {
      const [field] = changedFields;

      const [fieldName = '', entry = 0, keyName = ''] = field.name as Array<string | number>;

      if (fieldName === 'selectorList') onSelectorChange(JSON.stringify(selectorList));
      if (fieldName === 'assertionList') onChecksChange(JSON.stringify(form.getFieldValue('assertionList') || []));

      if (fieldName === 'assertionList' && keyName === 'key' && field.value) {
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
    [attributeList, form, selectorList]
  );

  return (
    <S.AssertionForm>
      <S.AssertionFormTitle strong>{isEditing ? 'Edit Assertion' : 'Add New Assertion'}</S.AssertionFormTitle>
      <Form
        name="assertion-form"
        form={form}
        initialValues={{
          remember: true,
          assertionList,
          selectorList,
        }}
        onFinish={onSubmit}
        autoComplete="off"
        layout="vertical"
        data-cy="assertion-form"
        onFieldsChange={onFieldsChange}
      >
        <div style={{marginBottom: 8}}>
          <Typography.Text>Filter to limit the span(s) included in this assertion</Typography.Text>
        </div>
        <S.SelectorInputContainer>
          <Form.Item
            name="selectorList"
            rules={[{required: true, message: 'At least one selector is required'}]}
            data-tour={GuidedTourService.getStep(GuidedTours.Assertion, Steps.Selectors)}
          >
            <AssertionFormSelectorInput attributeList={attributeList} />
          </Form.Item>
          <Form.Item name="pseudoSelector">
            <AssertionFormPseudoSelectorInput />
          </Form.Item>
        </S.SelectorInputContainer>
        <div style={{marginBottom: 8}}>
          <Typography.Text>Define the checks to run against each span selected</Typography.Text>
        </div>
        <div data-tour={GuidedTourService.getStep(GuidedTours.Assertion, Steps.Checks)}>
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
            Add
          </Button>
        </S.AssertionFromActions>
      </Form>
    </S.AssertionForm>
  );
};

export default AssertionForm;
