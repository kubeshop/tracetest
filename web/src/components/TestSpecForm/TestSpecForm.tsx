import {Button, Form, Input, Tag} from 'antd';
import React, {useState} from 'react';

import {ADVANCE_SELECTORS_DOCUMENTATION_URL} from 'constants/Common.constants';
import {CompareOperator} from 'constants/Operator.constants';
import {useAppSelector} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import SpanSelectors from 'selectors/Span.selectors';
import OperatorService from 'services/Operator.service';
import {TAssertion} from 'types/Assertion.types';
import {singularOrPlural} from 'utils/Common';
import AssertionCheckList from './AssertionCheckList';
import useAssertionFormValues from './hooks/useAssertionFormValues';
import useOnFieldsChange from './hooks/useOnFieldsChange';
import SelectorInput from './SelectorInput';
import * as S from './TestSpecForm.styled';

export interface IValues {
  assertionList?: TAssertion[];
  selector?: string;
}

interface IProps {
  defaultValues?: IValues;
  isEditing?: boolean;
  onCancel(): void;
  onSubmit(values: IValues): void;
  runId: string;
  testId: string;
}

const TestSpecForm = ({
  defaultValues: {
    assertionList = [
      {
        attribute: '',
        comparator: OperatorService.getOperatorSymbol(CompareOperator.EQUALS),
        expected: '',
      },
    ],
    selector = '',
  } = {},
  isEditing = false,
  onCancel,
  onSubmit,
  runId,
  testId,
}: IProps) => {
  const [form] = Form.useForm<IValues>();
  const [isValid, setIsValid] = useState(false);

  const spanIdList = useAppSelector(SpanSelectors.selectAffectedSpans);
  const attributeList = useAppSelector(state =>
    AssertionSelectors.selectAttributeList(state, testId, runId, spanIdList)
  );
  const {currentAssertionList} = useAssertionFormValues(form);

  const onFieldsChange = useOnFieldsChange({
    form,
    attributeList,
  });

  return (
    <S.AssertionForm>
      <S.AssertionFormHeader>
        <S.AssertionFormTitle>{isEditing ? 'Edit Test Spec' : 'Add Test Spec'}</S.AssertionFormTitle>
      </S.AssertionFormHeader>

      <Form<IValues>
        name="assertion-form"
        form={form}
        initialValues={{
          remember: true,
          assertionList,
          selector,
        }}
        onFinish={onSubmit}
        autoComplete="off"
        layout="vertical"
        data-cy="assertion-form"
        onFieldsChange={onFieldsChange}
      >
        <S.FormSection>
          <S.FormSectionRow1>
            <S.FormSectionTitle $noMargin>1. Select Spans</S.FormSectionTitle>
            <Tag color="blue">{`${spanIdList.length} ${singularOrPlural('span', spanIdList.length)} selected`}</Tag>
          </S.FormSectionRow1>
          <S.FormSectionRow>
            <S.FormSectionText>Specify which spans to assert using the </S.FormSectionText>
            <a href={ADVANCE_SELECTORS_DOCUMENTATION_URL} target="_blank">
              Selector Language
            </a>
          </S.FormSectionRow>
          <SelectorInput form={form} testId={testId} runId={runId} onValidSelector={setIsValid} />
        </S.FormSection>

        <S.FormSection>
          <S.FormSectionTitle>2. Add Assertions</S.FormSectionTitle>
          <S.FormSectionRow>
            <S.FormSectionText>Add assertions using the attributes from the selected spans</S.FormSectionText>
          </S.FormSectionRow>
          <Form.List name="assertionList">
            {(fields, {add, remove}) => (
              <AssertionCheckList
                assertionList={currentAssertionList}
                form={form}
                fields={fields}
                add={add}
                remove={remove}
                attributeList={attributeList}
              />
            )}
          </Form.List>
        </S.FormSection>

        <S.FormSection>
          <S.FormSectionTitle>3. Test Spec name (optional)</S.FormSectionTitle>
          <S.FormSectionRow>
            <S.FormSectionText>Identify your spec by giving it a name</S.FormSectionText>
          </S.FormSectionRow>
          <Form.Item name="name">
            <Input placeholder="Test Spec name" />
          </Form.Item>
        </S.FormSection>

        <S.AssertionFromActions>
          <Button onClick={onCancel}>Cancel</Button>
          <Button type="primary" disabled={!isValid} onClick={form.submit} data-cy="assertion-form-submit-button">
            Save Test Spec
          </Button>
        </S.AssertionFromActions>
      </Form>
    </S.AssertionForm>
  );
};

export default TestSpecForm;
