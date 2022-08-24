import {Button, Form, Typography} from 'antd';
import {ADVANCE_SELECTORS_DOCUMENTATION_URL} from 'constants/Common.constants';
import {CompareOperator} from 'constants/Operator.constants';
import React, {useState} from 'react';
import {useAppSelector} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import SpanSelectors from 'selectors/Span.selectors';
import OperatorService from 'services/Operator.service';
import {TAssertion} from 'types/Assertion.types';
import AffectedSpanControls from '../Diagram/components/DAG/AffectedSpanControls';
import {TooltipQuestion} from '../TooltipQuestion/TooltipQuestion';
import * as S from './AssertionForm.styled';
import AssertionFormCheckList from './AssertionFormCheckList';
import AssertionFormSelector from './AssertionFormSelector';
import useAssertionFormValues from './hooks/useAssertionFormValues';
import useOnFieldsChange from './hooks/useOnFieldsChange';

export interface IValues {
  assertionList?: TAssertion[];
  selector?: string;
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
    selector = '',
  } = {},
  onSubmit,
  onCancel,
  isEditing = false,
  testId,
  runId,
}) => {
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
        <S.AssertionFormTitle>{isEditing ? 'Edit Assertion' : 'Add New Assertion'}</S.AssertionFormTitle>
        <S.AffectedSpansContainer>
          <AffectedSpanControls />
          <S.AffectedSpansLabel>selected span(s)</S.AffectedSpansLabel>
        </S.AffectedSpansContainer>
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
        <S.AdvancedSelectorContainer>
          <Typography.Text>Filter to limit the span(s) included in this assertion</Typography.Text>
          <TooltipQuestion
            margin={0}
            title={`
            You can decide which spans will be tested by this assertion by altering the filter.
            Use the dropdown to the right to select the first matching span, last, n-th, or all.
            `}
          />
          <S.ReferenceLink>
            <a href={ADVANCE_SELECTORS_DOCUMENTATION_URL} target="_blank">
              Query Language Reference
            </a>
          </S.ReferenceLink>
        </S.AdvancedSelectorContainer>
        <AssertionFormSelector form={form} testId={testId} runId={runId} onValidSelector={setIsValid} />

        <div style={{marginBottom: 8}}>
          <Typography.Text>
            Define the checks to run as some checks are against the collection ({spanIdList.length} spans)
          </Typography.Text>
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
                form={form}
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
          <Button type="primary" disabled={!isValid} onClick={form.submit} data-cy="assertion-form-submit-button">
            {isEditing ? 'Save' : 'Add'}
          </Button>
        </S.AssertionFromActions>
      </Form>
    </S.AssertionForm>
  );
};

export default AssertionForm;
