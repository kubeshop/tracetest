import {PlusOutlined} from '@ant-design/icons';
import {FormInstance} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';
import {TAssertion} from '../../types/Assertion.types';
import {TSpanFlatAttribute} from '../../types/Span.types';
import {AssertionCheck} from './AssertionCheck';
import {useGetOTELSemanticConventionAttributesInfo} from './hooks/useGetOTELSemanticConventionAttributesInfo';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';

interface IProps {
  form: FormInstance<IValues>;
  fields: FormListFieldData[];
  assertionList: TAssertion[];

  add(): void;

  remove(name: number): void;

  attributeList: TSpanFlatAttribute[];
}

const AssertionCheckList: React.FC<IProps> = ({form, fields, add, remove, attributeList, assertionList}) => {
  const reference = useGetOTELSemanticConventionAttributesInfo();

  return (
    <S.AssertionsContainer>
      <S.CheckContainer>
        {fields.map(({key, name, ...field}, index) => {
          return (
            <AssertionCheck
              key={key}
              form={form}
              remove={remove}
              field={field}
              attributeList={attributeList}
              name={name}
              index={index}
              assertionList={assertionList}
              reference={reference}
            />
          );
        })}
      </S.CheckContainer>
      <S.AddCheckButton
        icon={<PlusOutlined />}
        onClick={() => {
          CreateAssertionModalAnalyticsService.onAddCheck();
          add();
        }}
        data-cy="add-assertion-form-add-check"
      >
        Add new
      </S.AddCheckButton>
    </S.AssertionsContainer>
  );
};

export default AssertionCheckList;
