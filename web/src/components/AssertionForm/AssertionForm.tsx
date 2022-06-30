import {Typography, Form, Button, Switch} from 'antd';
import React, {useState} from 'react';
import {CompareOperator} from 'constants/Operator.constants';
import {useAppSelector} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import OperatorService from 'services/Operator.service';
import {TAssertion, TPseudoSelector, TSpanSelector} from 'types/Assertion.types';
import {TooltipQuestion} from '../TooltipQuestion/TooltipQuestion';
import * as S from './AssertionForm.styled';
import AssertionFormCheckList from './AssertionFormCheckList';
import AssertionFormSelector from './AssertionFormSelector';
import SpanSelectors from '../../selectors/Span.selectors';
import useOnFieldsChange from './hooks/useOnFieldsChange';
import AffectedSpanControls from '../Diagram/components/DAG/AffectedSpanControls';

export interface IValues {
  assertionList?: TAssertion[];
  selectorList: TSpanSelector[];
  pseudoSelector?: TPseudoSelector;
  selector?: string;
  isAdvancedSelector?: boolean;
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
    selector = '',
    isAdvancedSelector = false,
  } = {},
  onSubmit,
  onCancel,
  isEditing = false,
  testId,
  runId,
}) => {
  const [form] = Form.useForm<IValues>();
  const currentAssertionList = Form.useWatch('assertionList', form) || [];
  const [isValid, setIsValid] = useState(false);

  const spanIdList = useAppSelector(SpanSelectors.selectAffectedSpans);
  const attributeList = useAppSelector(state =>
    AssertionSelectors.selectAttributeList(state, testId, runId, spanIdList)
  );

  const onFieldsChange = useOnFieldsChange({
    form,
    attributeList,
  });

  return (
    <S.AssertionForm>
      <S.AssertionFormHeader>
        <S.AssertionFormTitle strong>{isEditing ? 'Edit Assertion' : 'Add New Assertion'}</S.AssertionFormTitle>
        <S.AffectedSpansContainer>
          <AffectedSpanControls />
          <Typography.Text style={{color: '#61175e'}}>selected span(s)</Typography.Text>
        </S.AffectedSpansContainer>
      </S.AssertionFormHeader>
      <Form<IValues>
        name="assertion-form"
        form={form}
        initialValues={{
          remember: true,
          assertionList,
          selectorList,
          pseudoSelector,
          isAdvancedSelector,
          selector,
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
        <S.AdvancedSelectorContainer>
          <Typography.Text>Mode</Typography.Text>
          <Form.Item name="isAdvancedSelector" noStyle>
            <Switch
              checkedChildren="Advanced"
              unCheckedChildren="Wizard"
              disabled={isAdvancedSelector}
              defaultChecked={isAdvancedSelector}
              data-cy="mode-selector-switch"
            />
          </Form.Item>
          <TooltipQuestion
            margin={0}
            title={`
            You can decided if you want to use the wizard to create the span selector or the query language. 
            `}
          />
        </S.AdvancedSelectorContainer>
        <AssertionFormSelector
          selectorList={selectorList}
          pseudoSelector={pseudoSelector}
          form={form}
          testId={testId}
          runId={runId}
          isEditing={isEditing}
          onValidSelector={setIsValid}
        />

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
          <Button type="primary" disabled={!isValid} onClick={form.submit} data-cy="assertion-form-submit-button">
            {isEditing ? 'Save' : 'Add'}
          </Button>
        </S.AssertionFromActions>
      </Form>
    </S.AssertionForm>
  );
};

export default AssertionForm;
