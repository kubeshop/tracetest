// import CodeMirror from '@uiw/react-codemirror';
import {Form, FormInstance} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {capitalize} from 'lodash';
import {SupportedEditors} from '../../constants/Editor.constants';
import {CompareOperator} from '../../constants/Operator.constants';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import OperatorService from '../../services/Operator.service';
import {TAssertion} from '../../types/Assertion.types';
import {TSpanFlatAttribute} from '../../types/Span.types';
import Editor from '../Editor';
// import useEditorTheme from '../Editor/hooks/useEditorTheme';
import {OtelReference} from './hooks/useGetOTELSemanticConventionAttributesInfo';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';
// import {useExpectedInputLanguage} from './useExpectedInputLanguage';

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

export const AssertionCheck = ({field, index, name, remove}: IProps) => {
  return (
    <S.Container>
      <S.FieldsContainer>
        <Form.Item
          {...field}
          style={{margin: 0, width: 0, flexBasis: '30%', paddingLeft: 8}}
          name={[name, 'comparator']}
          rules={[{required: true, message: 'Operator is required'}]}
          data-cy="assertion-check-operator"
          initialValue={operatorList[0].value}
          noStyle
        >
          <Editor type={SupportedEditors.Expression} placeholder="Expression" />
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
