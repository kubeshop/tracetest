import {Form, FormInstance} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {SupportedEditors} from 'constants/Editor.constants';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import Editor from '../Editor';
import {OtelReference} from './hooks/useGetOTELSemanticConventionAttributesInfo';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';

interface IProps {
  remove(name: number): void;
  reference: OtelReference;
  form: FormInstance<IValues>;
  field: Pick<FormListFieldData, never>;
  name: number;
  attributeList: TSpanFlatAttribute[];
  index: number;
  assertions: string[];
}

export const AssertionCheck = ({field, index, name, remove}: IProps) => {
  return (
    <S.Container>
      <S.FieldsContainer>
        <Form.Item
          {...field}
          name={name}
          rules={[{required: true, message: 'Operator is required'}]}
          data-cy="assertion-input"
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
