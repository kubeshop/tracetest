import {Form, FormInstance, Select} from 'antd';
import {EditorView} from '@codemirror/view';
import {useCallback, useMemo, useState} from 'react';
import {delay} from 'lodash';
import {Completion, startCompletion} from '@codemirror/autocomplete';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {SupportedEditors} from 'constants/Editor.constants';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import OperatorService from 'services/Operator.service';
import {CompareOperator} from 'constants/Operator.constants';
import AssertionSelectors from 'selectors/Assertion.selectors';
import {useAppSelector} from 'redux/hooks';
import {Editor} from 'components/Inputs';
import {OtelReference} from './hooks/useGetOTELSemanticConventionAttributesInfo';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';
import AssertionCheckValue from './AssertionCheckValue';

interface IProps {
  remove(name: number): void;
  reference: OtelReference;
  form: FormInstance<IValues>;
  field: Pick<FormListFieldData, never>;
  name: number;
  attributeList: TSpanFlatAttribute[];
  runId: number;
  testId: string;
  spanIdList: string[];
}

const operatorList = Object.values(CompareOperator).map(value => ({
  value: OperatorService.getOperatorSymbol(value),
  label: OperatorService.getOperatorSymbol(value),
}));

export const AssertionCheck = ({field, name, remove, attributeList, runId, testId, spanIdList, form}: IProps) => {
  const [selectedAttributeKey, setSelectedAttributeKey] = useState('');
  const onAttributeFocus = useCallback((view: EditorView) => {
    if (!view?.state.doc.length) delay(() => startCompletion(view!), 0);
  }, []);

  const valueList = useAppSelector(state =>
    AssertionSelectors.selectAttributeValueList(state, testId, runId, spanIdList, selectedAttributeKey)
  );

  const onSelectAttribute = useCallback(
    ({label}: Completion) => {
      const attributeKey = attributeList.find(({key}) => key === label.replace('attr:', ''))?.key || '';

      setSelectedAttributeKey(attributeKey);
    },
    [attributeList]
  );

  const selector = Form.useWatch('selector') || '';
  const editorContext = useMemo(() => {
    return {
      runId,
      testId,
      selector,
    };
  }, [runId, selector, testId]);

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
            context={editorContext}
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
        <AssertionCheckValue
          form={form}
          valueList={valueList}
          editorContext={editorContext}
          name={name}
          field={field}
        />
      </S.FieldsContainer>
      <S.ActionContainer>
        <S.DeleteCheckIcon
          onClick={() => {
            CreateAssertionModalAnalyticsService.onRemoveCheck();
            remove(name);
          }}
        />
      </S.ActionContainer>
    </S.Container>
  );
};
