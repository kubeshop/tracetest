import CodeMirror from '@uiw/react-codemirror';
import {Form, FormInstance, Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {capitalize} from 'lodash';
import {CompareOperator} from '../../constants/Operator.constants';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import OperatorService from '../../services/Operator.service';
import {TAssertion} from '../../types/Assertion.types';
import {TSpanFlatAttribute} from '../../types/Span.types';
import useEditorTheme from '../AdvancedEditor/hooks/useEditorTheme';
import {ExpectedInputContainer} from './ExpectedInputContainer';
import {AttributeField} from './Fields/AttributeField';
import {OtelReference} from './hooks/useGetOTELSemanticConventionAttributesInfo';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';
import {useExpectedInputLanguage} from './useExpectedInputLanguage';

const operatorList = Object.values(CompareOperator).map(value => ({
  value: OperatorService.getOperatorSymbol(value),
  label: capitalize(OperatorService.getOperatorName(value)),
}));

interface IProps {
  remove(name: number): void;
  reference: OtelReference;
  form: FormInstance<IValues>;
  field: Pick<FormListFieldData, never>;
  name: number;
  attributeList: TSpanFlatAttribute[];
  index: number;
  assertions: TAssertion[];
}

export const AssertionCheck = ({attributeList, field, index, name, assertions, form, remove, reference}: IProps) => {
  const extensionList = useExpectedInputLanguage();
  const editorTheme = useEditorTheme();
  return (
    <S.Container>
      <S.FieldsContainer>
        <AttributeField field={field} name={name} attributeList={attributeList} reference={reference} />
        <Form.Item
          {...field}
          style={{margin: 0, width: 0, flexBasis: '30%', paddingLeft: 8}}
          name={[name, 'comparator']}
          rules={[{required: true, message: 'Operator is required'}]}
          data-cy="assertion-check-operator"
          initialValue={operatorList[0].value}
        >
          <S.Select style={{margin: 0}} placeholder="Assertion Type">
            {operatorList.map(({value, label}) => (
              <Select.Option key={value} value={value}>
                {label}
              </Select.Option>
            ))}
          </S.Select>
        </Form.Item>

        <ExpectedInputContainer>
          <Form.Item
            {...field}
            style={{margin: 0}}
            rules={[{required: true, message: 'Value is required'}]}
            shouldUpdate
          >
            {value => {
              const assertionValues = value.getFieldValue(`assertions`);
              return (
                <CodeMirror
                  id="assertion-check-value"
                  basicSetup={{lineNumbers: false}}
                  data-cy="assertion-check-value"
                  value={assertionValues[name]?.expected}
                  extensions={extensionList}
                  onChange={val =>
                    form.setFieldsValue({
                      assertions: assertions.map((a, i) => (i === name ? {...a, expected: val} : a)),
                    })
                  }
                  spellCheck={false}
                  theme={editorTheme}
                  placeholder="Expected Value"
                  style={{width: '100%'}}
                />
              );
            }}
          </Form.Item>
        </ExpectedInputContainer>
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
