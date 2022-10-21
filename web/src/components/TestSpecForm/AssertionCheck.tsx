import {Form, FormInstance, Select} from 'antd';
import {EditorView} from '@codemirror/view';
import {useCallback} from 'react';
import {delay} from 'lodash';
import {Completion, startCompletion} from '@codemirror/autocomplete';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {SupportedEditors} from 'constants/Editor.constants';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import {TStructuredAssertion} from 'types/Assertion.types';
import OperatorService from 'services/Operator.service';
import {CompareOperator} from 'constants/Operator.constants';
import Editor from '../Editor';
import {OtelReference} from './hooks/useGetOTELSemanticConventionAttributesInfo';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';
import AssertionService from '../../services/Assertion.service';

interface IProps {
  remove(name: number): void;
  reference: OtelReference;
  form: FormInstance<IValues>;
  field: Pick<FormListFieldData, never>;
  name: number;
  attributeList: TSpanFlatAttribute[];
  index: number;
  assertions: TStructuredAssertion[];
}

const operatorList = Object.values(CompareOperator).map(value => ({
  value: OperatorService.getOperatorSymbol(value),
  label: OperatorService.getOperatorSymbol(value),
}));

export const AssertionCheck = ({field, index, name, remove, form, assertions, attributeList}: IProps) => {
  const onAttributeFocus = useCallback((view: EditorView) => {
    if (!view?.state.doc.length) delay(() => startCompletion(view!), 0);
  }, []);

  const onSelectAttribute = useCallback(
    ({label}: Completion) => {
      const value = attributeList.find(({key}) => key === label.replace('attr:', ''))?.value || '';

      form.setFieldsValue({
        assertions: assertions.map((assertion, i) =>
          i === name ? {...assertion, right: AssertionService.extractExpectedString(value) || ''} : assertion
        ),
      });
    },
    [assertions, attributeList, form, name]
  );

  return (
    <S.Container>
      <S.FieldsContainer>
        <Form.Item
          {...field}
          name={[name, 'left']}
          rules={[{required: true, message: 'An attribute is required'}]}
          style={{margin: 0}}
          data-cy="assertion-check-attribute"
        >
          <Editor
            type={SupportedEditors.Expression}
            placeholder="Attribute"
            onFocus={onAttributeFocus}
            onSelectAutocompleteOption={onSelectAttribute}
          />
        </Form.Item>
        <Form.Item
          {...field}
          name={[name, 'comparator']}
          rules={[{required: true, message: 'Operator is required'}]}
          style={{margin: 0}}
          initialValue={operatorList[0].value}
        >
          <S.Select data-cy="assertion-check-operator" style={{margin: 0}} placeholder="Assertion Type">
            {operatorList.map(({value, label}) => (
              <Select.Option key={value} value={value}>
                {label}
              </Select.Option>
            ))}
          </S.Select>
        </Form.Item>
        <Form.Item
          {...field}
          name={[name, 'right']}
          rules={[{required: true, message: 'Expected value is required'}]}
          data-cy="assertion-check-value"
          style={{margin: 0}}
        >
          <Editor type={SupportedEditors.Expression} placeholder="Expected Value" />
        </Form.Item>
      </S.FieldsContainer>
      <S.ActionContainer>
        {index !== 0 && (
          <S.DeleteCheckIcon
            onClick={() => {
              CreateAssertionModalAnalyticsService.onRemoveCheck();
              remove(name);
            }}
          />
        )}
      </S.ActionContainer>
    </S.Container>
  );
};
